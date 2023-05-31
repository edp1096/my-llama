#include <cstring>

#include "common.h"
#include "myllama.h"

void* init_container() {
    myllama_container* c = new myllama_container;

    c->gptparams = new gpt_params;
    c->ctxparams = new llama_context_params(llama_context_default_params());
    c->session_tokens = new std::vector<llama_token>;

    return c;
}

// void allocate_tokens(void* container, char* text, bool add_bos) {
void allocate_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;
    // std::vector<llama_token> tokens(strlen(text) + (int)add_bos);
    std::vector<llama_token> tokens(((llama_context_params*)c->ctxparams)->n_ctx);

    c->tokens = (void*)tokens.data();
    c->n_tokens = tokens.size();
}

int* get_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;

    return (int*)c->tokens;
}

// void prepare_candidates(void* container, int n_vocab) {
int prepare_candidates(void* container, int n_vocab) {
    myllama_container* c = (myllama_container*)container;

    std::vector<llama_token_data> candidates;
    candidates.reserve(n_vocab);
    for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
        candidates.emplace_back(llama_token_data{token_id, c->logits[token_id], 0.0f});
    }

    llama_token_data_array candidates_da = {candidates.data(), candidates.size(), false};

    int id = 0;

    if ((llama_context*)c->ctx != nullptr) {
        auto next_token = llama_sample_token((llama_context*)c->ctx, &candidates_da);
        id = (int)next_token;
    }

    return id;
}