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

void allocate_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;
    gpt_params* gptparams = (gpt_params*)c->gptparams;

    std::vector<llama_token> tokens(gptparams->n_ctx);

    c->tokens = (void*)tokens.data();
    c->n_tokens = tokens.size();
}

int* get_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;

    return (int*)c->tokens;
}

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
