#include <cstring>
#include <algorithm>

#include "common.h"
#include "myllama.h"

void* init_container() {
    myllama_container* c = new myllama_container;

    c->gptparams = (void*)new gpt_params;
    c->ctxparams = (void*)new llama_context_params(llama_context_default_params());
    c->session_tokens = (void*)new std::vector<llama_token>;

    return c;
}

// void allocate_tokens(void* container, char* text, bool add_bos) {
void allocate_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;
    gpt_params* gptparams = (gpt_params*)c->gptparams;

    // std::vector<llama_token> tokens(strlen(text) + (int)add_bos);
    std::vector<llama_token> tokens(gptparams->n_ctx);

    c->tokens = (void*)tokens.data();
    c->n_tokens = tokens.size();
}

int* get_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;

    return (int*)c->tokens;
}

// // int prepare_candidates(void* container, int n_vocab) {
// void prepare_candidates(void* container, int n_vocab) {
//     myllama_container* c = (myllama_container*)container;
//     llama_context* ctx = (llama_context*)c->ctx;

//     // float* logits = c->logits;
//     float* logits = llama_get_logits(ctx);
//     n_vocab = llama_n_vocab(ctx);

//     std::vector<llama_token_data> candidates;
//     candidates.reserve(n_vocab);
//     for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
//         candidates.emplace_back(llama_token_data{token_id, logits[token_id], 0.0f});
//     }

//     llama_token_data_array candidates_p = {candidates.data(), candidates.size(), false};

//     c->candidates = (void*)new llama_token_data_array(candidates_p);
// }

void prepare_candidates(void* container, int n_vocab) {
    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    float* logits = c->logits;

    llama_token_data_array* candidates_p = (llama_token_data_array*)malloc(sizeof(llama_token_data_array));
    candidates_p->data = (llama_token_data*)malloc(sizeof(llama_token_data) * n_vocab);
    candidates_p->size = n_vocab;
    candidates_p->sorted = false;

    for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
        candidates_p->data[token_id] = llama_token_data{token_id, logits[token_id], 0.0f};
    }

    c->candidates = (void*)candidates_p;
}

void mini_run_main(void* container, int n_past, char* prompt) {
    // const char* model_fname = "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin";
    // const char* prompt = "The quick brown fox";

    // auto gptparams = new gpt_params;

    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;
    gpt_params* gptparams = (gpt_params*)c->gptparams;
    llama_context_params* ctxparams = (llama_context_params*)c->ctxparams;

    gptparams->seed = 42;
    gptparams->n_threads = 4;
    gptparams->repeat_last_n = 64;
    gptparams->n_predict = 16;

    // llama_context_params* ctxparams = new llama_context_params(llama_context_default_params());

    ctxparams->n_ctx = gptparams->n_ctx;
    ctxparams->seed = gptparams->seed;
    ctxparams->f16_kv = gptparams->memory_f16;
    ctxparams->use_mmap = gptparams->use_mmap;
    ctxparams->use_mlock = gptparams->use_mlock;

    // auto n_past = 0;

    // // init
    // llama_context* ctx = llama_init_from_file(model_fname, *ctxparams);

    auto tokens = std::vector<llama_token>(gptparams->n_ctx);
    auto n_prompt_tokens = llama_tokenize(ctx, prompt, tokens.data(), tokens.size(), true);

    if (n_prompt_tokens < 1) {
        fprintf(stderr, "%s : failed to tokenize prompt\n", __func__);
        return;
    }

    // print tokens.data() numbers loop
    for (auto i = 0; i < n_prompt_tokens; i++) {
        printf("%d ", tokens.data()[i]);
    }

    // evaluate prompt
    llama_eval(ctx, tokens.data(), n_prompt_tokens, n_past, gptparams->n_threads);
    n_past += n_prompt_tokens;

    printf("tokens.size(): %d\n", tokens.size());

    for (auto i = 0; i < gptparams->n_predict; i++) {
        auto logits = llama_get_logits(ctx);
        auto n_vocab = llama_n_vocab(ctx);
        // printf("logits[0], n_predict: %f, %d\n", logits[0], gptparams->n_predict);

        std::vector<llama_token_data> candidates;
        candidates.reserve(n_vocab);

        for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
            candidates.emplace_back(llama_token_data{token_id, logits[token_id], 0.0f});
        }
        llama_token_data_array candidates_p = {candidates.data(), candidates.size(), false};

        auto next_token = llama_sample_token(ctx, &candidates_p);
        auto next_token_str = llama_token_to_str(ctx, next_token);

        printf("%s", next_token_str);
        // printf("n_predict, n_vocab, n_tokens, n_past: %d, %d, %d, %d\n", gptparams->n_predict, n_vocab, next_token, n_past);
        if (llama_eval(ctx, &next_token, 1, n_past, gptparams->n_threads)) {
            fprintf(stderr, "\n%s : failed to evaluate\n", __func__);
            return;
        }
        n_past += 1;
    }

    printf("\n\n");
}