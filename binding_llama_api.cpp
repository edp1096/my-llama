#include <stdlib.h>
#include <string.h>

#include "llama.h"
#include "binding.h"
#include "binding_llama_api.h"

// Not done
void llama_api_context_default_params() {
    llama_context_default_params();
}

bool llama_api_mmap_supported() {
    return llama_mmap_supported();
}

bool llama_api_mlock_supported() {
    return llama_mlock_supported();
}

void llama_api_init_backend() {
    llama_init_backend();
}

int64_t llama_api_time_us() {
    return llama_time_us();
}

// Not done
void llama_api_init_from_file() {
    llama_init_from_file("", {});
}

void llama_api_free(void* container) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_free((llama_context*)c->ctx);
    }
}

// Not done
int llama_api_model_quantize() {
    return llama_model_quantize("", "", {}, 0);
}

// Not done
int llama_api_get_kv_cache_token_count() {
    return llama_get_kv_cache_token_count({});
}

void llama_api_set_rng_seed(void* container, int seed) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_set_rng_seed((llama_context*)c->ctx, seed);
    }
}

// Not done
int llama_api_get_state_size() {
    return llama_get_state_size({});
}

// Not done
int llama_api_copy_state_data() {
    return llama_copy_state_data({}, {});
}

// Not done
int llama_api_set_state_data() {
    return llama_set_state_data({}, 0);
}

// Not done
bool llama_api_load_session_file() {
    return llama_load_session_file({}, "", {}, 0, 0);
}

// Not done
bool llama_api_save_session_file() {
    return llama_save_session_file({}, "", {}, 0);
}

// Not done
int llama_api_eval() {
    return llama_eval({}, {}, 0, 0, 0);
}

// Not done
int llama_api_tokenize() {
    return llama_tokenize({}, "", {}, 0, false);
}

int llama_api_n_vocab(void* container) {
    int n_vocab = 0;
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        n_vocab = llama_n_vocab((llama_context*)c->ctx);
    }

    return n_vocab;
}

int llama_api_n_ctx(void* container) {
    int n_ctx = 0;
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        n_ctx = llama_n_ctx((llama_context*)c->ctx);
    }

    return n_ctx;
}

int llama_api_n_embd(void* container) {
    int n_embd = 0;
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        n_embd = llama_n_embd((llama_context*)c->ctx);
    }

    return n_embd;
}

// Not done
float* llama_api_get_logits(void* container) {
    float* logits;
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        logits = llama_get_logits((llama_context*)c->ctx);
    }

    return logits;
}

// Not done
float* llama_api_get_embeddings() {
    return llama_get_embeddings({});
}

// Not done
char* llama_api_token_to_str() {
    const char* result = llama_token_to_str({}, 0);

    char* c_result = (char*)malloc(strlen(result) + 1);
    strcpy(c_result, result);

    return c_result;
}

int llama_api_token_bos() {
    return (int)llama_token_bos();
}

int llama_api_token_eos() {
    return (int)llama_token_eos();
}

int llama_api_token_nl() {
    return (int)llama_token_nl();
}

// Not done
void llama_api_sample_repetition_penalty() {
    llama_sample_repetition_penalty({}, {}, {}, 0, 0);
}

// Not done
void llama_api_sample_frequency_and_presence_penalties() {
    llama_sample_frequency_and_presence_penalties({}, {}, {}, 0, 0, 0);
}

// Not done
void llama_api_sample_softmax() {
    llama_sample_softmax({}, {});
}

// Not done
void llama_api_sample_top_k() {
    llama_sample_top_k({}, {}, 0, 0);
}

// Not done
void llama_api_sample_top_p() {
    llama_sample_top_p({}, {}, 0, 0);
}

// Not done
void llama_api_sample_tail_free() {
    llama_sample_tail_free({}, {}, 0, 0);
}

// Not done
void llama_api_sample_typical() {
    llama_sample_typical({}, {}, 0, 0);
}

// Not done
void llama_api_sample_temperature() {
    llama_sample_temperature({}, {}, 0);
}

// Not done
void llama_api_sample_token_mirostat_v2() {
    llama_sample_token_mirostat_v2({}, {}, 0, 0, 0);
}

// Not done
void llama_api_sample_token_greedy() {
    llama_sample_token_greedy({}, {});
}

// Not done
void llama_api_sample_token() {
    llama_sample_token({}, {});
}

void llama_api_print_timings(void* container) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_print_timings((llama_context*)c->ctx);
    }
}

void llama_api_reset_timings(void* container) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_reset_timings((llama_context*)c->ctx);
    }
}

char* llama_api_print_system_info() {
    const char* result = llama_print_system_info();

    char* c_result = (char*)malloc(strlen(result) + 1);
    strcpy(c_result, result);

    return c_result;
}
