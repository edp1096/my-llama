#include <cstring>

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

// void prepare_candidates(void* container, int n_vocab) {
int prepare_candidates(void* container, int n_vocab_no) {
    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    // float* logits = c->logits;
    float* logits = llama_get_logits(ctx);
    int n_vocab = llama_n_vocab(ctx);

    std::vector<llama_token_data> candidates;
    candidates.reserve(n_vocab);
    for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
        candidates.emplace_back(llama_token_data{token_id, logits[token_id], 0.0f});
    }

    // llama_token_data_array candidates_da = {candidates.data(), candidates.size(), false};
    llama_token_data_array candidates_p = {candidates.data(), candidates.size(), false};
    auto next_token = llama_sample_token(ctx, &candidates_p);

    int id = 0;
    id = (int)next_token;

    // if (ctx != nullptr) {
    //     auto next_token = llama_sample_token(ctx, &candidates_p);
    //     id = (int)next_token;
    // }

    return id;
}