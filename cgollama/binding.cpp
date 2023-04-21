#include "common.h"
#include "llama.h"
#include "ggml.h"
#include "binding.h"

#include <cstring>
#include <iostream>

struct llama_kv_cache {
    struct ggml_tensor* k;
    struct ggml_tensor* v;

    int n;  // number of tokens currently in the cache
};

struct llama_model {
    struct llama_kv_cache kv_self;
};

struct llama_context {
    llama_model model;
};

void* llama_init_container() {
    variables_container* c = new variables_container;
    c->params = new gpt_params;

    return c;
}

llama_context* llama_init_context(gpt_params* params) {
    auto lparams = llama_context_default_params();

    lparams.n_ctx = params->n_ctx;
    lparams.n_parts = params->n_parts;
    lparams.seed = params->seed;
    lparams.f16_kv = params->memory_f16;
    lparams.use_mmap = params->use_mmap;
    lparams.use_mlock = params->use_mlock;

    llama_context* ctx = llama_init_from_file(params->model.c_str(), lparams);

    return ctx;
}

bool llama_load_model(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;
    gpt_params* params = (gpt_params*)c->params;

    if (params->seed <= 0) {
        params->seed = time(NULL);
    }

    std::mt19937 rng(params->seed);
    if (params->random_prompt) {
        params->prompt = gpt_random_prompt(rng);
    }

    llama_context* ctx = llama_init_context(params);
    if (ctx == nullptr) {
        fprintf(stderr, "%s : failed to load model\n", __func__);
        return result;
    }

    c->ctx = ctx;

    result = true;
    return result;
}

void llama_save_kv_dump_experiment(void* container) {
    variables_container* c = (variables_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    // Save ctx->model.kv_self.k->data and ctx->model.kv_self.v->data and ctx->model.kv_self.n to file
    FILE* fp_write = fopen("dump_kv.bin", "wb");
    fwrite(ctx->model.kv_self.k->data, ggml_nbytes(ctx->model.kv_self.k), 1, fp_write);
    fwrite(ctx->model.kv_self.v->data, ggml_nbytes(ctx->model.kv_self.v), 1, fp_write);
    fwrite(&ctx->model.kv_self.n, sizeof(ctx->model.kv_self.n), 1, fp_write);
    fclose(fp_write);
}

void llama_load_kv_dump_experiment(void* container) {
    variables_container* c = (variables_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    // Load ctx->model.kv_self.k->data and ctx->model.kv_self.v->data and ctx->model.kv_self.n from file
    FILE* fp_read = fopen("dump_kv.bin", "rb");
    fread(ctx->model.kv_self.k->data, ggml_nbytes(ctx->model.kv_self.k), 1, fp_read);
    fread(ctx->model.kv_self.v->data, ggml_nbytes(ctx->model.kv_self.v), 1, fp_read);
    fread(&ctx->model.kv_self.n, sizeof(ctx->model.kv_self.n), 1, fp_read);
    fclose(fp_read);
}

bool llama_make_ready_to_predict(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;
    gpt_params* params = (gpt_params*)c->params;

    c->is_interacting = false;
    c->embd = new std::vector<llama_token>;

    params->prompt.insert(0, 1, ' '); // Add a space in front of the first character to match OG llama tokenizer behavior

    c->embd_inp = new std::vector<llama_token>(::llama_tokenize((llama_context*)c->ctx, params->prompt, true)); // tokenize the prompt
    c->n_ctx = llama_n_ctx((llama_context*)c->ctx);

    if ((int)((std::vector<llama_token>*)c->embd_inp)->size() > c->n_ctx - 4) {
        fprintf(stderr, "%s: error: prompt is too long (%d tokens, max %d)\n", __func__, (int)((std::vector<llama_token>*)c->embd_inp)->size(), c->n_ctx - 4);
        return result;
    }

    // number of tokens to keep when resetting context
    if (params->n_keep < 0 || params->n_keep >(int)((std::vector<llama_token>*)c->embd_inp)->size() || params->instruct) {
        params->n_keep = (int)((std::vector<llama_token>*)c->embd_inp)->size();
    }

    // enable interactive mode if reverse prompt or interactive start is specified
    if (params->antiprompt.size() != 0 || params->interactive_start) {
        params->interactive = true;
    }

    // determine newline token
    c->llama_token_newline = new std::vector<llama_token>(::llama_tokenize((llama_context*)c->ctx, "\n", false));

    // TODO: replace with ring-buffer
    c->last_n_tokens = new std::vector<llama_token>(c->n_ctx);
    std::fill(((std::vector<llama_token>*)c->last_n_tokens)->begin(), ((std::vector<llama_token>*)c->last_n_tokens)->end(), 0);

    c->input_noecho = false;

    c->n_past = 0;
    c->n_remain = params->n_predict;
    c->n_consumed = 0;


    result = true;
    return result;
}

void llama_init_params(void* container) {
    gpt_params* params = (gpt_params*)((variables_container*)container)->params;

    params->interactive = true;
    params->interactive_start = params->interactive;
    params->antiprompt = {};

    params->n_threads = 6;
    params->n_predict = 512;
    params->use_mlock = true;

    // params->prompt = "The quick brown fox jumps over the lazy dog.";
    // params->n_predict = 100;
    // params->n_keep = 0;
    // params->instruct = false;
    // params->interactive_start = false;
    // params->temp = 1.0;
    // params->top_k = 40;
    // params->top_p = 0.0;
    // params->seed = 0;
}

bool llama_predict_tokens(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;

    std::vector<llama_token>* last_n_tokens = (std::vector<llama_token>*)c->last_n_tokens;
    std::vector<llama_token>* llama_token_newline = (std::vector<llama_token>*)c->llama_token_newline;
    std::vector<llama_token>* embd = (std::vector<llama_token>*)c->embd;
    std::vector<llama_token>* embd_inp = (std::vector<llama_token>*)c->embd_inp;
    llama_context* ctx = (llama_context*)c->ctx;
    gpt_params* params = (gpt_params*)c->params;

    // predict
    if (embd->size() > 0) {
        if (c->n_past + (int)embd->size() > c->n_ctx) {
            const int n_left = c->n_past - params->n_keep;

            c->n_past = params->n_keep;

            // insert n_left/2 tokens at the start of embd from last_n_tokens
            embd->insert(embd->begin(), last_n_tokens->begin() + c->n_ctx - n_left / 2 - embd->size(), last_n_tokens->end() - embd->size());
        }

        if (llama_eval(ctx, embd->data(), embd->size(), c->n_past, params->n_threads)) {
            fprintf(stderr, "%s : failed to eval\n", __func__);
            return result;
        }
    }

    c->n_past += embd->size();
    embd->clear();

    if ((int)embd_inp->size() <= c->n_consumed) {
        // out of user input, sample next token
        const int32_t top_k = params->top_k;
        const float   top_p = params->top_p;
        const float   temp = params->temp;
        const float   repeat_penalty = params->repeat_penalty;

        llama_token id = 0;

        {
            auto logits = llama_get_logits(ctx);

            if (params->ignore_eos) {
                logits[llama_token_eos()] = 0;
            }

            id = llama_sample_top_p_top_k(ctx,
                last_n_tokens->data() + c->n_ctx - params->repeat_last_n,
                params->repeat_last_n, top_k, top_p, temp, repeat_penalty);

            last_n_tokens->erase(last_n_tokens->begin());
            last_n_tokens->push_back(id);
        }

        // replace end of text token with newline token when in interactive mode
        if (id == llama_token_eos() && params->interactive && !params->instruct) {
            id = llama_token_newline->front();
            if (params->antiprompt.size() != 0) {
                // tokenize and inject first reverse prompt
                const auto first_antiprompt = ::llama_tokenize(ctx, params->antiprompt.front(), false);
                embd_inp->insert(embd_inp->end(), first_antiprompt.begin(), first_antiprompt.end());
            }
        }

        embd->push_back(id); // add it to the context

        c->input_noecho = false; // echo this to console

        --c->n_remain; // decrement remaining sampling budget
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

bool llama_receive_input(void* container) {
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
            return result; // input stream is bad or EOF received
        }
        win32_utf8_encode(wline, line);

        if (line.empty() || line.back() != '\\') {
            another_line = false;
        } else {
            line.pop_back(); // Remove the continue character
        }
        buffer += line + '\n'; // Append the line to the result
    } while (another_line);
    // buffer = "Do you know kimchi?";

    // Add tokens to embd only if the input buffer is non-empty. Entering a empty line lets the user pass control back
    if (buffer.length() > 1) {
        auto line_inp = ::llama_tokenize(ctx, buffer, false);
        embd_inp->insert(embd_inp->end(), line_inp.begin(), line_inp.end());

        c->n_remain -= line_inp.size();
    }

    c->input_noecho = true; // do not echo this again

    result = true;
    return result;
}

bool llama_append_input(void* container) {
    bool result = false;
    variables_container* c = (variables_container*)container;

    std::vector<llama_token>* embd_inp = (std::vector<llama_token>*)c->embd_inp;
    llama_context* ctx = (llama_context*)c->ctx;

    char* buffer = (char*)c->user_input;

    // Add tokens to embd only if the input buffer is non-empty. Entering a empty line lets the user pass control back
    if (strlen(buffer) > 1) {
        std::vector<llama_token> line_inp = ::llama_tokenize(ctx, buffer, false);
        embd_inp->insert(embd_inp->end(), line_inp.begin(), line_inp.end());

        c->n_remain -= line_inp.size();
    }


    c->input_noecho = true; // do not echo this again

    result = true;
    return result;
}

bool llama_wait_or_continue(void* container) {
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
        // // bool result = llama_receive_input(c);
        // bool result = llama_append_input(c);
        // if (!result) {
        //     return false;
        // }

        // // printf("%s", c->user_input);
        // c->user_input = NULL;
    }

    if (c->n_past > 0) {
        c->is_interacting = false;
    }

    return true;
}

int llama_get_embed_id(void* container, int index) {
    variables_container* c = (variables_container*)container;
    return ((std::vector<llama_token>*)c->embd)->at(index);
}

char* llama_get_embed_string(void* container, int id) {
    variables_container* c = (variables_container*)container;
    return const_cast<char*>(llama_token_to_str((llama_context*)c->ctx, id));
}


void llama_free_params(void* container) {
    gpt_params* params = (gpt_params*)((variables_container*)container)->params;
    delete params;
}

void llama_free_model(void* container) {
    llama_context* ctx = (llama_context*)((variables_container*)container)->ctx;
    llama_free(ctx);
}


/* Getters */
int llama_get_n_remain(void* container) {
    return ((variables_container*)container)->n_remain;
}

int llama_get_params_n_predict(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->n_predict;
}

bool llama_get_noecho(void* container) {
    return ((variables_container*)container)->input_noecho;
}

int llama_get_embd_size(void* container) {
    return (int)((std::vector<llama_token>*)((variables_container*)container)->embd)->size();
}

int llama_get_embd_inp_size(void* container) {
    return (int)((std::vector<llama_token>*)((variables_container*)container)->embd_inp)->size();
}

int llama_get_n_consumed(void* container) {
    return ((variables_container*)container)->n_consumed;
}

bool llama_get_params_interactive_start(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->interactive_start;
}

bool llama_get_params_interactive(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->interactive;
}

/* Getters - gpt_params */
int llama_get_params_n_threads(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->n_threads;
}

/* Getters - gpt_params / sampling parameters */
int llama_get_params_top_k(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->top_k;
}

float llama_get_params_top_p(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->top_p;
}

float llama_get_params_temper(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->temp;
}

float llama_get_params_repeat_penalty(void* container) {
    return ((gpt_params*)((variables_container*)container)->params)->repeat_penalty;
}


/* Setters */
void llama_set_params_interactive_start(void* container) {
    bool interactive = ((gpt_params*)((variables_container*)container)->params)->interactive;
    ((gpt_params*)((variables_container*)container)->params)->interactive_start = interactive;
}

void llama_set_is_interacting(void* container, bool is_interacting) {
    ((variables_container*)container)->is_interacting = is_interacting;
}

void llama_set_n_remain(void* container, int n_predict) {
    ((variables_container*)container)->n_remain = n_predict;
}

void llama_set_model_path(void* container, char* path) {
    ((gpt_params*)((variables_container*)container)->params)->model = path;
}

void llama_set_params_antiprompt(void* container, char* antiprompt) {
    ((gpt_params*)((variables_container*)container)->params)->antiprompt.push_back(strdup(antiprompt));
}

void llama_set_params_prompt(void* container, char* prompt) {
    ((gpt_params*)((variables_container*)container)->params)->prompt = strdup(prompt);
}

void llama_set_user_input(void* container, const char* user_input) {
    ((variables_container*)container)->user_input = strdup(user_input);
}


/* Setters - gpt_params */
void llama_set_params_n_threads(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->n_threads = value;
}

/* Setters - gpt_params / sampling parameters */
void llama_set_params_top_k(void* container, int value) {
    ((gpt_params*)((variables_container*)container)->params)->top_k = value;
}

void llama_set_params_top_p(void* container, float value) {
    ((gpt_params*)((variables_container*)container)->params)->top_p = value;
}

void llama_set_params_temper(void* container, float value) {
    ((gpt_params*)((variables_container*)container)->params)->temp = value;
}

void llama_set_params_repeat_penalty(void* container, float value) {
    ((gpt_params*)((variables_container*)container)->params)->repeat_penalty = value;
}


bool llama_check_prompt_or_continue(void* container) {
    bool result = true;
    variables_container* c = (variables_container*)container;

    // in interactive mode, and not currently processing queued inputs. check if we should prompt the user for more
    if (((gpt_params*)c->params)->interactive && (int)((std::vector<llama_token>*)c->embd_inp)->size() <= c->n_consumed) {
        result = llama_wait_or_continue(c);
    }

    if (result) {
        // end of text token
        if (((std::vector<llama_token>*)c->embd)->back() == llama_token_eos()) {
            c->is_interacting = false;
            result = false;

            printf("DONE?");
        }
    }

    return result;
}

void llama_dropback_user_input(void* container) {
    variables_container* c = (variables_container*)container;

    // In interactive mode, respect the maximum number of tokens and drop back to user input when reached.
    if (c->n_remain <= 0 && ((gpt_params*)c->params)->n_predict != -1) {
        c->n_remain = ((gpt_params*)c->params)->n_predict;
    }
}
