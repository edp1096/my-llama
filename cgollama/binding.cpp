#include "common.h"
#include "llama.h"
#include "binding.h"

#include <cassert>
#include <cinttypes>
#include <cmath>
#include <cstdio>
#include <cstring>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

#if defined (__unix__) || (defined (__APPLE__) && defined (__MACH__))
#include <signal.h>
#include <unistd.h>
#elif defined (_WIN32)
#include <signal.h>
#endif

#if defined (__unix__) || (defined (__APPLE__) && defined (__MACH__)) || defined (_WIN32)
void sigint_handler(int signo) {
    if (signo == SIGINT) {
        _exit(130);
    }
}
#endif

int llama_get_remain_count(void* pred_vars_v) {
    pred_variables* pred_vars_p = (pred_variables*)pred_vars_v;

    return (int)pred_vars_p->n_remain;
}

int llama_get_id(void* pred_vars_ptr, int index) {
    pred_variables* pred_vars_p = (pred_variables*)pred_vars_ptr;

    std::vector<llama_token>* embedding = (std::vector<llama_token>*)(pred_vars_p->embd);
    std::vector<llama_token> embd = *embedding;

    return (int)embd[index];
}

char* llama_get_embed_string(void* pred_vars_ptr, int id) {
    pred_variables* pred_vars_p = (pred_variables*)pred_vars_ptr;
    llama_context* ctx = (llama_context*)(pred_vars_p->context);

    std::string predicted = llama_token_to_str(ctx, id);

    return strdup(predicted.c_str());
}

bool llama_check_token_end(void* pred_vars_ptr) {
    pred_variables* pred_vars_p = (pred_variables*)pred_vars_ptr;

    std::vector<llama_token>* embedding = (std::vector<llama_token>*)(pred_vars_p->embd);
    std::vector<llama_token> embd = *embedding;

    bool result = false;

    // end of text token
    if (embd.back() == llama_token_eos()) {
        result = true;
    }

    return result;
}

int llama_get_embedding_ids(void* params_ptr, void* pred_vars_ptr) {
    gpt_params* params_p = (gpt_params*)params_ptr;
    gpt_params params = *params_p;

    pred_variables* pred_vars_p = (pred_variables*)pred_vars_ptr;

    int n_past = pred_vars_p->n_past;
    int n_remain = pred_vars_p->n_remain;
    int n_consumed = pred_vars_p->n_consumed;
    int n_ctx = pred_vars_p->n_ctx;

    std::vector<llama_token>* embedding = (std::vector<llama_token>*)(pred_vars_p->embd);
    std::vector<llama_token> embd = *embedding;

    std::vector<llama_token>* embedding_inp = (std::vector<llama_token>*)(pred_vars_p->embedding_inp);
    std::vector<llama_token> embd_inp = *embedding_inp;

    std::vector<llama_token>* last_n_tokens_p = (std::vector<llama_token>*)(pred_vars_p->last_n_tokens);
    std::vector<llama_token> last_n_tokens = *last_n_tokens_p;

    llama_context* ctx = (llama_context*)(pred_vars_p->context);

    // determine newline token
    auto llama_token_newline = ::llama_tokenize(ctx, "\n", false);

    // predict
    if (embd.size() > 0) {
        // infinite text generation via context swapping
        // if we run out of context:
        // - take the n_keep first tokens from the original prompt (via n_past)
        // - take half of the last (n_ctx - n_keep) tokens and recompute the logits in a batch
        if (n_past + (int)embd.size() > n_ctx) {
            const int n_left = n_past - params.n_keep;
            n_past = params.n_keep;

            // insert n_left/2 tokens at the start of embd from last_n_tokens
            embd.insert(embd.begin(), last_n_tokens.begin() + n_ctx - n_left / 2 - embd.size(), last_n_tokens.end() - embd.size());
        }

        if (llama_eval(ctx, embd.data(), embd.size(), n_past, params.n_threads)) {
            fprintf(stderr, "%s : failed to eval\n", __func__);
            return 0;
        }
    }

    n_past += embd.size();
    embd.clear();

    if ((int)embd_inp.size() <= n_consumed) {
        // out of user input, sample next token
        const int32_t top_k = params.top_k;
        const float   top_p = params.top_p;
        const float   temp = params.temp;
        const float   repeat_penalty = params.repeat_penalty;

        llama_token id = 0;

        {
            auto logits = llama_get_logits(ctx);

            if (params.ignore_eos) {
                logits[llama_token_eos()] = 0;
            }

            id = llama_sample_top_p_top_k(ctx,
                last_n_tokens.data() + n_ctx - params.repeat_last_n,
                params.repeat_last_n, top_k, top_p, temp, repeat_penalty);

            last_n_tokens.erase(last_n_tokens.begin());
            last_n_tokens.push_back(id);
        }

        // replace end of text token with newline token when in interactive mode
        if (id == llama_token_eos()) {
            id = llama_token_newline.front();
            if (params.antiprompt.size() != 0) {
                // tokenize and inject first reverse prompt
                const auto first_antiprompt = ::llama_tokenize(ctx, params.antiprompt.front(), false);
                embd_inp.insert(embd_inp.end(), first_antiprompt.begin(), first_antiprompt.end());
            }
        }

        embd.push_back(id); // add it to the context

        --n_remain; // decrement remaining sampling budget
    } else {
        // some user input remains from prompt or interaction, forward it to processing
        while ((int)embd_inp.size() > n_consumed) {
            embd.push_back(embd_inp[n_consumed]);
            last_n_tokens.erase(last_n_tokens.begin());
            last_n_tokens.push_back(embd_inp[n_consumed]);
            ++n_consumed;
            if ((int)embd.size() >= params.n_batch) {
                break;
            }
        }
    }

    pred_vars_p->n_remain = n_remain;
    pred_vars_p->n_past = n_past;
    pred_vars_p->n_consumed = n_consumed;

    *embedding = embd;
    pred_vars_p->embd = embedding;

    *last_n_tokens_p = last_n_tokens;
    pred_vars_p->last_n_tokens = last_n_tokens_p;

    std::vector<int> ids = (std::vector<int>)(embd);

    return ids.size();
}

void llama_default_signal_action() {
#if defined (_WIN32)
    signal(SIGINT, SIG_DFL);
#endif
}

void* llama_prepare_pred_vars(void* params_ptr, void* state_pr) {
    gpt_params* params_p = (gpt_params*)params_ptr;
    llama_context* ctx = (llama_context*)state_pr;

    gpt_params params = *params_p;

    if (params.seed <= 0) {
        params.seed = time(NULL);
    }

    std::mt19937 rng(params.seed);

    params.prompt.insert(0, 1, ' '); // Add a space in front of the first character to match OG llama tokenizer behavior

    auto embd_inp = ::llama_tokenize(ctx, params.prompt, true); // tokenize the prompt

    const int n_ctx = llama_n_ctx(ctx);

    // number of tokens to keep when resetting context
    if (params.n_keep < 0 || params.n_keep >(int)embd_inp.size() || params.instruct) {
        params.n_keep = (int)embd_inp.size();
    }

    std::vector<llama_token>* last_n_tokens_p = new std::vector<llama_token>(n_ctx);
    std::fill(last_n_tokens_p->begin(), last_n_tokens_p->end(), 0);

    pred_variables* pred_vars = new pred_variables;

    pred_vars->n_past = 0;
    pred_vars->n_remain = params.n_predict;
    pred_vars->n_consumed = 0;

    std::vector<llama_token>* embd = new std::vector<llama_token>;
    pred_vars->embd = embd;

    std::vector<llama_token>* embd_inp_ptr = new std::vector<llama_token>;
    *embd_inp_ptr = embd_inp;
    pred_vars->embedding_inp = embd_inp_ptr;

    pred_vars->context = ctx;
    pred_vars->last_n_tokens = last_n_tokens_p;
    pred_vars->n_ctx = n_ctx;

    return pred_vars;
}

void llama_free_model(void* state_ptr) {
    llama_context* ctx = (llama_context*)state_ptr;
    llama_free(ctx);
}

void llama_free_params(void* params_ptr) {
    gpt_params* params = (gpt_params*)params_ptr;
    delete params;
}


void* llama_allocate_params(const char* prompt, const char* antiprompt, int seed, int threads, int tokens, int top_k,
    float top_p, float temp, float repeat_penalty, int repeat_last_n, bool ignore_eos, bool memory_f16) {

    gpt_params* params = new gpt_params;
    params->seed = seed;
    params->n_batch = 512;
    params->n_threads = threads;
    params->n_predict = tokens;
    params->repeat_last_n = repeat_last_n;

    params->top_k = top_k;
    params->top_p = top_p;
    params->memory_f16 = memory_f16;
    params->temp = temp;
    params->repeat_penalty = repeat_penalty;

    params->prompt = prompt;
    params->ignore_eos = ignore_eos;

    // Add anti-prompt
    params->antiprompt = std::vector<std::string>();
    if (antiprompt != NULL) {
        std::string antiprompt_str = antiprompt;
        std::string delimiter = "|";

        size_t pos = 0;
        std::string token;
        while ((pos = antiprompt_str.find(delimiter)) != std::string::npos) {
            token = antiprompt_str.substr(0, pos);
            params->antiprompt.push_back(token);
            antiprompt_str.erase(0, pos + delimiter.length());
        }
        params->antiprompt.push_back(antiprompt_str);
    }

    return params;
}

void* llama_update_params(void* params_ptr, const char* prompt) {
    gpt_params* params = (gpt_params*)params_ptr;

    params->prompt = prompt;

    return params;
}

void* load_model(const char* fname, int n_ctx, int n_parts, int n_seed, bool memory_f16, bool mlock) {
    // load the model
    auto lparams = llama_context_default_params();

    lparams.n_ctx = n_ctx;
    lparams.n_parts = n_parts;
    lparams.seed = n_seed;
    lparams.f16_kv = memory_f16;
    lparams.use_mlock = mlock;

    return llama_init_from_file(fname, lparams);
}
