#include <iostream>
#include <cassert>
#include <cstring>

// #include "llama.h"
#include "common.h"
#include "binding.h"

// Windows not yet support so, override it
#ifdef _WIN32
int32_t get_num_physical_cores() {
    return 1;
}
#endif

// TODO: not great allocating this every time
std::vector<llama_token> binding_tokenize(struct llama_context* ctx, const std::string& text, bool add_bos) {
    // initialize to prompt numer of chars, since n_tokens <= n_prompt_chars
    std::vector<llama_token> res(text.size() + (int)add_bos);
    int n = llama_tokenize(ctx, text.c_str(), res.data(), res.size(), add_bos);
    assert(n >= 0);
    res.resize(n);

    return res;
}

void* bd_init_container() {
    variables_container* c = new variables_container;
    c->params = new gpt_params;
    c->session_tokens = new std::vector<llama_token>;

    return c;
}

struct llama_context* llama_init_from_gpt_params(const gpt_params& params) {
    auto lparams = llama_context_default_params();

    lparams.n_ctx = params.n_ctx;
    // lparams.n_parts = params.n_parts;
#ifndef USE_OLD_GGML
    lparams.n_gpu_layers = params.n_gpu_layers;
#endif
    lparams.seed = params.seed;
    lparams.f16_kv = params.memory_f16;
    lparams.use_mmap = params.use_mmap;
    lparams.use_mlock = params.use_mlock;
    lparams.logits_all = params.perplexity;
    lparams.embedding = params.embedding;

    llama_context* lctx = llama_init_from_file(params.model.c_str(), lparams);

    if (lctx == NULL) {
        fprintf(stderr, "%s: error: failed to load model '%s'\n", __func__, params.model.c_str());
        return NULL;
    }

    if (!params.lora_adapter.empty()) {
        int err = llama_apply_lora_from_file(lctx,
                                             params.lora_adapter.c_str(),
                                             params.lora_base.empty() ? NULL : params.lora_base.c_str(),
                                             params.n_threads);
        if (err != 0) {
            fprintf(stderr, "%s: error: failed to apply lora adapter\n", __func__);
            return NULL;
        }
    }

    return lctx;
}

bool bd_load_model(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;
    gpt_params* params = (gpt_params*)c->params;

    if (params->seed < 0) {
        params->seed = time(NULL);
    }

#ifndef USE_OLD_GGML
    params->n_gpu_layers = 33;  // for 3060ti, CUDA only
    printf("n_gpu_layers: %d\n", params->n_gpu_layers);
#endif

    // llama_context* ctx = binding_init_context(params);
    llama_context* ctx = llama_init_from_gpt_params(*params);
    if (ctx == nullptr) {
        fprintf(stderr, "%s : failed to load model\n", __func__);
        return result;
    }

    c->ctx = ctx;

    result = true;
    return result;
}

bool bd_allocate_variables(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;
    gpt_params* params = (gpt_params*)c->params;

    c->is_interacting = false;
    c->embd = new std::vector<llama_token>;

    params->prompt.insert(0, 1, ' ');  // Add a space in front of the first character to match OG llama tokenizer behavior

    c->embd_inp = new std::vector<llama_token>(::binding_tokenize((llama_context*)c->ctx, params->prompt, true));  // tokenize the prompt
    c->n_ctx = llama_n_ctx((llama_context*)c->ctx);

    if ((int)((std::vector<llama_token>*)c->embd_inp)->size() > c->n_ctx - 4) {
        fprintf(stderr, "%s: error: prompt is too long (%d tokens, max %d)\n", __func__, (int)((std::vector<llama_token>*)c->embd_inp)->size(), c->n_ctx - 4);
        return result;
    }

    // number of tokens to keep when resetting context
    if (params->n_keep < 0 || params->n_keep > (int)((std::vector<llama_token>*)c->embd_inp)->size() || params->instruct) {
        params->n_keep = (int)((std::vector<llama_token>*)c->embd_inp)->size();
    }

    // enable interactive mode if reverse prompt or interactive start is specified
    if (params->antiprompt.size() != 0 || params->interactive_first) {
        params->interactive = true;
    }

    // determine newline token
    c->llama_token_newline = new std::vector<llama_token>(::binding_tokenize((llama_context*)c->ctx, "\n", false));

    // TODO: replace with ring-buffer
    c->last_n_tokens = new std::vector<llama_token>(c->n_ctx);
    std::fill(((std::vector<llama_token>*)c->last_n_tokens)->begin(), ((std::vector<llama_token>*)c->last_n_tokens)->end(), 0);

    c->input_noecho = false;

    c->n_past = 0;
    c->n_remain = params->n_predict;
    c->n_consumed = 0;
    c->n_session_consumed = 0;

    result = true;
    return result;
}

void bd_init_params(void* container) {
    gpt_params* params = (gpt_params*)((variables_container*)container)->params;

    params->interactive = true;
    params->interactive_first = params->interactive;
    params->antiprompt = {};

    // params->n_predict = 512;
    params->use_mmap = false;
    params->use_mlock = true;
}

bool bd_predict_tokens(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;

    std::vector<llama_token>* last_n_tokens = (std::vector<llama_token>*)c->last_n_tokens;
    std::vector<llama_token>* llama_token_newline = (std::vector<llama_token>*)c->llama_token_newline;
    std::vector<llama_token>* embd = (std::vector<llama_token>*)c->embd;
    std::vector<llama_token>* embd_inp = (std::vector<llama_token>*)c->embd_inp;
    llama_context* ctx = (llama_context*)c->ctx;
    gpt_params* params = (gpt_params*)c->params;

    std::vector<llama_token>* session_tokens = (std::vector<llama_token>*)c->session_tokens;

    // predict
    if (embd->size() > 0) {
        if (c->n_past + (int)embd->size() > c->n_ctx) {
            const int n_left = c->n_past - params->n_keep;

            // c->n_past = params->n_keep;
            // always keep the first token - BOS
            c->n_past = std::max(1, params->n_keep);

            // insert n_left/2 tokens at the start of embd from last_n_tokens
            embd->insert(embd->begin(), last_n_tokens->begin() + c->n_ctx - n_left / 2 - embd->size(), last_n_tokens->end() - embd->size());
        }

        // if (llama_eval(ctx, embd->data(), embd->size(), c->n_past, params->n_threads)) {
        //     fprintf(stderr, "%s : failed to eval\n", __func__);
        //     return result;
        // }

        // try to reuse a matching prefix from the loaded session instead of re-eval (via n_past)
        if (c->n_session_consumed < (int)session_tokens->size()) {
            size_t i = 0;
            for (; i < embd->size(); i++) {
                if (embd[i] != session_tokens[c->n_session_consumed]) {
                    session_tokens->resize(c->n_session_consumed);
                    break;
                }

                c->n_past++;
                c->n_session_consumed++;

                if (c->n_session_consumed >= (int)session_tokens->size()) {
                    ++i;
                    break;
                }
            }

            if (i > 0) {
                embd->erase(embd->begin(), embd->begin() + i);
            }
        }

        // evaluate tokens in batches
        // embd is typically prepared beforehand to fit within a batch, but not always
        for (int i = 0; i < (int)embd->size(); i += params->n_batch) {
            int n_eval = (int)embd->size() - i;
            if (n_eval > params->n_batch) {
                n_eval = params->n_batch;
            }
            if (llama_eval(ctx, embd->data() + i, n_eval, c->n_past, params->n_threads)) {
                fprintf(stderr, "%s : failed to eval\n", __func__);
                return result;
            }
            c->n_past += n_eval;
        }

        // if (embd.size() > 0 && !path_session.empty()) {
        if (embd->size() > 0) {
            session_tokens->insert(session_tokens->end(), embd->begin(), embd->end());
            c->n_session_consumed = session_tokens->size();
        }
    }

    // c->n_past += embd->size();
    embd->clear();

    if ((int)embd_inp->size() <= c->n_consumed) {
        // out of user input, sample next token
        const int32_t top_k = params->top_k;
        const float top_p = params->top_p;
        const float tfs_z = params->tfs_z;
        const float temp = params->temp;
        const float typical_p = params->typical_p;
        const int32_t repeat_last_n = params->repeat_last_n < 0 ? c->n_ctx : params->repeat_last_n;
        const float repeat_penalty = params->repeat_penalty;
        const float alpha_presence = params->presence_penalty;
        const float alpha_frequency = params->frequency_penalty;
        const int mirostat = params->mirostat;
        const float mirostat_tau = params->mirostat_tau;
        const float mirostat_eta = params->mirostat_eta;
        const bool penalize_nl = params->penalize_nl;

        llama_token id = 0;

        {
            auto logits = llama_get_logits(ctx);
            auto n_vocab = llama_n_vocab(ctx);

            if (params->penalize_nl) {
                params->logit_bias[llama_token_eos()] = -INFINITY;
            }

            std::vector<llama_token_data> candidates;
            candidates.reserve(n_vocab);
            for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
                candidates.emplace_back(llama_token_data{token_id, logits[token_id], 0.0f});
            }

            llama_token_data_array candidates_p = {candidates.data(), candidates.size(), false};

            // Apply penalties
            float nl_logit = logits[llama_token_nl()];
            auto last_n_repeat = std::min(std::min((int)last_n_tokens->size(), repeat_last_n), c->n_ctx);
            llama_sample_repetition_penalty(ctx, &candidates_p,
                                            last_n_tokens->data() + last_n_tokens->size() - last_n_repeat,
                                            last_n_repeat, repeat_penalty);
            llama_sample_frequency_and_presence_penalties(ctx, &candidates_p,
                                                          last_n_tokens->data() + last_n_tokens->size() - last_n_repeat,
                                                          last_n_repeat, alpha_frequency, alpha_presence);
            if (!penalize_nl) {
                logits[llama_token_nl()] = nl_logit;
            }

            if (temp <= 0) {
                // Greedy sampling
                id = llama_sample_token_greedy(ctx, &candidates_p);
            } else {
                if (mirostat == 1) {
                    static float mirostat_mu = 2.0f * mirostat_tau;
                    const int mirostat_m = 100;
                    llama_sample_temperature(ctx, &candidates_p, temp);
                    id = llama_sample_token_mirostat(ctx, &candidates_p, mirostat_tau, mirostat_eta, mirostat_m, &mirostat_mu);
                } else if (mirostat == 2) {
                    static float mirostat_mu = 2.0f * mirostat_tau;
                    llama_sample_temperature(ctx, &candidates_p, temp);
                    id = llama_sample_token_mirostat_v2(ctx, &candidates_p, mirostat_tau, mirostat_eta, &mirostat_mu);
                } else {
                    // Temperature sampling
                    llama_sample_top_k(ctx, &candidates_p, top_k, 1);
                    llama_sample_tail_free(ctx, &candidates_p, tfs_z, 1);
                    llama_sample_typical(ctx, &candidates_p, typical_p, 1);
                    llama_sample_top_p(ctx, &candidates_p, top_p, 1);
                    llama_sample_temperature(ctx, &candidates_p, temp);
                    id = llama_sample_token(ctx, &candidates_p);
                }
            }

            last_n_tokens->erase(last_n_tokens->begin());
            last_n_tokens->push_back(id);
        }

        // replace end of text token with newline token when in interactive mode
        if (id == llama_token_eos() && params->interactive && !params->instruct) {
            id = llama_token_newline->front();
            if (params->antiprompt.size() != 0) {
                // tokenize and inject first reverse prompt
                const auto first_antiprompt = ::binding_tokenize(ctx, params->antiprompt.front(), false);
                embd_inp->insert(embd_inp->end(), first_antiprompt.begin(), first_antiprompt.end());
            }
        }

        embd->push_back(id);  // add it to the context

        c->input_noecho = false;  // echo this to console

        --c->n_remain;  // decrement remaining sampling budget
    } else {
        // some user input remains from prompt or interaction, forward it to processing
        while ((int)embd_inp->size() > c->n_consumed) {
            embd->push_back((*embd_inp)[c->n_consumed]);
            last_n_tokens->erase(last_n_tokens->begin());
            last_n_tokens->push_back((*embd_inp)[c->n_consumed]);
            ++c->n_consumed;
            if ((int)embd->size() >= params->n_batch) {
                break;
            }
        }
    }

    result = true;
    return result;
}

bool bd_receive_input(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;

    std::vector<llama_token>* embd_inp = (std::vector<llama_token>*)c->embd_inp;
    llama_context* ctx = (llama_context*)c->ctx;
    gpt_params* params = (gpt_params*)c->params;

    std::string buffer;
    if (!params->input_prefix.empty()) {
        buffer += params->input_prefix;
        printf("%s", buffer.c_str());
    }

    std::string line;
    bool another_line = true;
    do {
        std::wstring wline;
        if (!std::getline(std::wcin, wline)) {
            result = false;
            return result;  // input stream is bad or EOF received
        }
        // win32_utf8_encode(wline, line);

        if (line.empty() || line.back() != '\\') {
            another_line = false;
        } else {
            line.pop_back();  // Remove the continue character
        }
        buffer += line + '\n';  // Append the line to the result
    } while (another_line);
    // buffer = "What is the largest city in Earth?";

    // Add tokens to embd only if the input buffer is non-empty. Entering a empty line lets the user pass control back
    if (buffer.length() > 1) {
        auto line_inp = ::binding_tokenize(ctx, buffer, false);
        embd_inp->insert(embd_inp->end(), line_inp.begin(), line_inp.end());

        c->n_remain -= line_inp.size();
    }

    c->input_noecho = true;  // do not echo this again

    result = true;
    return result;
}

bool bd_append_input(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;

    std::vector<llama_token>* embd_inp = (std::vector<llama_token>*)c->embd_inp;
    llama_context* ctx = (llama_context*)c->ctx;

    char* buffer = (char*)c->user_input;

    // Add tokens to embd only if the input buffer is non-empty. Entering a empty line lets the user pass control back
    if (strlen(buffer) > 1) {
        std::vector<llama_token> line_inp = ::binding_tokenize(ctx, buffer, false);
        embd_inp->insert(embd_inp->end(), line_inp.begin(), line_inp.end());

        c->n_remain -= line_inp.size();
    }

    c->input_noecho = true;  // do not echo this again

    result = true;
    return result;
}

bool bd_wait_or_continue(void* container) {
    variables_container* c = (variables_container*)container;
    gpt_params* params = (gpt_params*)c->params;

    // check for reverse prompt
    if (params->antiprompt.size()) {
        std::string last_output;
        for (auto id : *(std::vector<llama_token>*)c->last_n_tokens) {
            last_output += llama_token_to_str((llama_context*)c->ctx, id);
        }

        // Check if each of the reverse prompts appears at the end of the output.
        for (std::string& antiprompt : params->antiprompt) {
            if (last_output.find(antiprompt.c_str(), last_output.length() - antiprompt.length(), antiprompt.length()) != std::string::npos) {
                c->is_interacting = true;
                fflush(stdout);
                break;
            }
        }
    }

    // Receive user input
    if (c->n_past > 0 && c->is_interacting) {
        return false;
    }

    if (c->n_past > 0) {
        c->is_interacting = false;
    }

    return true;
}

int bd_get_embed_id(void* container, int index) {
    variables_container* c = (variables_container*)container;
    return ((std::vector<llama_token>*)c->embd)->at(index);
}

char* bd_get_embed_string(void* container, int id) {
    variables_container* c = (variables_container*)container;
    return const_cast<char*>(llama_token_to_str((llama_context*)c->ctx, id));
}

/* Frees */
void bd_free_params(void* container) {
    gpt_params* params = (gpt_params*)((variables_container*)container)->params;
    if (params != NULL) {
        delete params;
    }
}

void bd_free_model(void* container) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_free((llama_context*)c->ctx);
    }
}

/* Getters */
int bd_get_n_remain(void* container) {
    return ((variables_container*)container)->n_remain;
}

int bd_get_params_n_predict(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->n_predict;
}

bool bd_get_noecho(void* container) {
    return ((variables_container*)container)->input_noecho;
}

int bd_get_embd_size(void* container) {
    return (int)((std::vector<llama_token>*)((variables_container*)container)->embd)->size();
}

int bd_get_embd_inp_size(void* container) {
    return (int)((std::vector<llama_token>*)((variables_container*)container)->embd_inp)->size();
}

int bd_get_n_consumed(void* container) {
    return ((variables_container*)container)->n_consumed;
}

bool bd_get_params_interactive_first(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->interactive_first;
}

bool bd_get_params_interactive(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->interactive;
}

/* Getters - gpt_params */
int bd_get_params_n_threads(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->n_threads;
}

/* Getters - gpt_params / sampling parameters */
int bd_get_params_top_k(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->top_k;
}

float bd_get_params_top_p(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->top_p;
}

float bd_get_params_temper(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->temp;
}

float bd_get_params_repeat_penalty(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->repeat_penalty;
}

/* Setters */
void bd_set_params_interactive_first(void* container) {
    bool interactive = ((gpt_params*)((variables_container*)container)->params)->interactive;
    ((gpt_params*)((variables_container*)container)->params)->interactive_first = interactive;
}

void bd_set_is_interacting(void* container, bool is_interacting) {
    ((variables_container*)container)->is_interacting = is_interacting;
}

void bd_set_n_remain(void* container, int n_predict) {
    ((variables_container*)container)->n_remain = n_predict;
}

void bd_set_model_path(void* container, char* path) {
    ((gpt_params*)((variables_container*)container)->params)->model = path;
}

void bd_set_params_antiprompt(void* container, char* antiprompt) {
    ((gpt_params*)((variables_container*)container)->params)->antiprompt.push_back(strdup(antiprompt));
}

void bd_set_params_prompt(void* container, char* prompt) {
    ((gpt_params*)((variables_container*)container)->params)->prompt = strdup(prompt);
}

void bd_set_user_input(void* container, const char* user_input) {
    ((variables_container*)container)->user_input = strdup(user_input);
}

/* Setters - gpt_params */
void bd_set_params_n_threads(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->n_threads = value;
}

void bd_set_params_use_mlock(void* container, bool value) {
    ((gpt_params*)((variables_container*)container)->params)->use_mlock = value;
}

/* Setters - gpt_params / sampling parameters */
void bd_set_params_n_ctx(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->n_ctx = value;
}

void bd_set_params_n_batch(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->n_batch = value;
}

void bd_set_sampling_method(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->mirostat = value;
}

void bd_set_params_top_k(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->top_k = value;
}

void bd_set_params_top_p(void* container, float value) {
    ((gpt_params*)((variables_container*)container)->params)->top_p = value;
}

void bd_set_params_temper(void* container, float value) {
    ((gpt_params*)((variables_container*)container)->params)->temp = value;
}

void bd_set_params_repeat_penalty(void* container, float value) {
    ((gpt_params*)((variables_container*)container)->params)->repeat_penalty = value;
}

/* State dump */
void bd_save_state(void* container, char* fname) {
    variables_container* c = (variables_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    size_t ctx_size = llama_get_state_size(ctx);
    uint8_t* state_mem = new uint8_t[ctx_size];
    llama_copy_state_data(ctx, state_mem);

    FILE* fp_write = fopen(fname, "wb");
    fwrite(&ctx_size, 1, sizeof(size_t), fp_write);
    fwrite(state_mem, 1, ctx_size, fp_write);
    fwrite(&c->last_n_tokens, 1, sizeof(int), fp_write);
    fwrite(&c->n_past, 1, sizeof(int), fp_write);
    fclose(fp_write);

    delete[] state_mem;
}

void bd_load_state(void* container, char* fname) {
    variables_container* c = (variables_container*)container;
    // llama_context* ctx = (llama_context*)c->ctx;

    FILE* fp_read = fopen(fname, "rb");
    size_t ctx_size = llama_get_state_size((llama_context*)c->ctx);

    size_t state_size;
    uint8_t* state_mem = new uint8_t[ctx_size];
    fread(&state_size, 1, sizeof(size_t), fp_read);

    if (state_size != ctx_size) {
        printf("Error: state size mismatch. Expected %zu, got %zu\n", ctx_size, state_size);
    }

    fread(state_mem, 1, state_size, fp_read);
    fread(&c->last_n_tokens, 1, sizeof(int), fp_read);
    fread(&c->n_past, 1, sizeof(int), fp_read);
    fclose(fp_read);

    llama_set_state_data((llama_context*)c->ctx, state_mem);
    delete[] state_mem;
}

void bd_save_session(void* container, char* fname) {
    variables_container* c = (variables_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    llama_save_session_file(ctx, fname, ((std::vector<llama_token>*)c->session_tokens)->data(), ((std::vector<llama_token>*)c->session_tokens)->size());
}

void bd_load_session(void* container, char* fname) {
    variables_container* c = (variables_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    // fopen to check for existing session
    FILE* fp = std::fopen(fname, "rb");
    if (fp != NULL) {
        std::fclose(fp);

        ((std::vector<llama_token>*)c->session_tokens)->resize(((gpt_params*)c->params)->n_ctx);
        size_t n_token_count_out = 0;
        if (!llama_load_session_file(ctx, fname, ((std::vector<llama_token>*)c->session_tokens)->data(), ((std::vector<llama_token>*)c->session_tokens)->capacity(), &n_token_count_out)) {
            fprintf(stderr, "%s: error: failed to load session file '%s'\n", __func__, fname);
            return;
        }
        ((std::vector<llama_token>*)c->session_tokens)->resize(n_token_count_out);

        fprintf(stderr, "%s: loaded a session with prompt size of %d tokens\n", __func__, (int)((std::vector<llama_token>*)c->session_tokens)->size());
    } else {
        fprintf(stderr, "%s: session file does not exist, will create\n", __func__);
    }
}

/* Others */
bool bd_check_prompt_or_continue(void* container) {
    bool result = true;
    variables_container* c = (variables_container*)container;

    // in interactive mode, and not currently processing queued inputs. check if we should prompt the user for more
    if (((gpt_params*)c->params)->interactive && (int)((std::vector<llama_token>*)c->embd_inp)->size() <= c->n_consumed) {
        result = bd_wait_or_continue(c);
    }

    if (result) {
        // end of text token
        if (((std::vector<llama_token>*)c->embd)->back() == llama_token_eos()) {
            c->is_interacting = false;
            result = false;
        }
    }

    return result;
}

void bd_dropback_user_input(void* container) {
    variables_container* c = (variables_container*)container;

    // In interactive mode, respect the maximum number of tokens and drop back to user input when reached.
    if (c->n_remain <= 0 && ((gpt_params*)c->params)->n_predict != -1) {
        c->n_remain = ((gpt_params*)c->params)->n_predict;
    }
}

void bd_print_timings(void* container) {
    variables_container* c = (variables_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    llama_print_timings(ctx);
}
