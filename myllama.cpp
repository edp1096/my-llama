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

void allocate_tokens(void* container, char* text, bool add_bos) {
    myllama_container* c = (myllama_container*)container;
    std::vector<llama_token> tokens(strlen(text) + (int)add_bos);

    c->tokens = (void*)tokens.data();
    c->n_tokens = (void*)new int(tokens.size());
}

int* get_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;

    return (int*)c->tokens;
}
