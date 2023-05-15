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
    lparams.n_parts = params.n_parts;
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
    params->n_gpu_layers = 32;  // for 3060ti, CUDA only
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


/* Setters */
void bd_set_model_path(void* container, char* path) {
    ((gpt_params*)((variables_container*)container)->params)->model = path;
}
