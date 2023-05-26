#include <stdlib.h>
#include <string.h>

// #include "llama.h"
#include "common.h"
#include "binding.h"  // including struct variables_container
#include "binding_llama_api.h"

void* llama_api_context_default_params() {
    return (void*)&llama_context_default_params();
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

void* llama_api_init_from_file(char* path_model, void* params_p) {
    llama_context_params params = *(llama_context_params*)params_p;

    return (void*)llama_init_from_file(path_model, params);
}

void llama_api_free(void* container) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_free((llama_context*)c->ctx);
    }
}

int llama_api_model_quantize(char* fname_inp, char* fname_out, int ftype_int, int nthread) {
    llama_ftype ftype = (llama_ftype)ftype_int;
    return llama_model_quantize(fname_inp, fname_out, ftype, nthread);
}

int llama_api_get_kv_cache_token_count(void* container) {
    int result = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        result = llama_get_kv_cache_token_count((llama_context*)c->ctx);
    }

    return result;
}

void llama_api_set_rng_seed(void* container, int seed) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_set_rng_seed((llama_context*)c->ctx, seed);
    }
}

int llama_api_get_state_size(void* container) {
    int result = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        result = (int)llama_get_state_size((llama_context*)c->ctx);
    }

    return result;        
}

int llama_api_copy_state_data(void* container, void* dst_p) {
    int result = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        result = (int)llama_copy_state_data((llama_context*)c->ctx, (uint8_t*)dst_p);
    }

    return result;
}

int llama_api_set_state_data(void* container, void* src_p) {
    int result = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        result = (int)llama_set_state_data((llama_context*)c->ctx, (uint8_t*)src_p);
    }

    return result;
}

bool llama_api_load_session_file(void* container, char* path_session, void* tokens_out, int n_token_capacity, int* n_token_count_out) {
    bool result = false;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        result = llama_load_session_file((llama_context*)c->ctx, path_session, (llama_token*)tokens_out, (size_t)n_token_capacity, (size_t*)n_token_count_out);
    }

    return result;
}

bool llama_api_save_session_file(void* container, char* path_session, void* tokens, int n_token_count) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        return llama_save_session_file((llama_context*)c->ctx, path_session, (llama_token*)tokens, (size_t)n_token_count);
    }
}

int llama_api_eval(void* container, int* tokens, int n_tokens, int n_past, int n_threads) {
    int result = 1;  // 0: success, 1: fail

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        result = llama_eval((llama_context*)c->ctx, tokens, n_tokens, n_past, (int32_t)n_threads);
    }

    return result;
}

int llama_api_tokenize(void* container, char* text, bool add_bos) {
    int result = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        std::vector<llama_token> tokens(strlen(text) + (int)add_bos);
        return llama_tokenize((llama_context*)c->ctx, text, tokens.data(), tokens.size(), add_bos);
    }

    return 0;
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

void* llama_api_get_logits(void* container) {
    void* logits;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        void* logits = (void*)llama_get_logits((llama_context*)c->ctx);
    }

    return logits;
}

void* llama_api_get_embeddings(void* container) {
    void* embeddings;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        void* embeddings = (void*)llama_get_embeddings((llama_context*)c->ctx);
    }

    return embeddings;
}

char* llama_api_token_to_str(void* container, int token) {
    llama_token token_id = (llama_token)token;
    char* c_result;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        const char* result = llama_token_to_str((llama_context*)c->ctx, token_id);
        c_result = (char*)malloc(strlen(result) + 1);
        strcpy(c_result, result);
    }

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

void llama_api_sample_repetition_penalty(void* container, void* candidates_a_p, void* last_tokens, int last_tokens_size, float penalty) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_repetition_penalty((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, (llama_token*)last_tokens, (size_t)last_tokens_size, penalty);
    }
}

void llama_api_sample_frequency_and_presence_penalties(void* container, void* candidates_a_p, void* last_tokens, int last_tokens_size, float alpha_frequency, float alpha_presence) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_frequency_and_presence_penalties((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, (llama_token*)last_tokens, (size_t)last_tokens_size, alpha_frequency, alpha_presence);
    }

}

void llama_api_sample_softmax(void* container, void* candidates_a_p) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_softmax((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p);
    }
}

void llama_api_sample_top_k(void* container, void* candidates_a_p, int top_k) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_top_k((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, top_k, 1);
    }
}

void llama_api_sample_top_p(void* container, void* candidates_a_p, float top_p) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_top_p((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, top_p, 1);
    }
}

void llama_api_sample_tail_free(void* container, void* candidates_a_p, float tfs_z) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_tail_free((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, tfs_z, 1);
    }
}

void llama_api_sample_typical(void* container, void* candidates_a_p, float typical_p) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_typical((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, typical_p, 1);
    }
}

void llama_api_sample_temperature(void* container, void* candidates_a_p, float temperature) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_temperature((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, temperature);
    }
}

void llama_api_sample_token_mirostat_v2(void* container, void* candidates_a_p, float mirostat_tau, float mirostat_eta, float mirostat_mu) {
    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        llama_sample_token_mirostat_v2((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p, mirostat_tau, mirostat_eta, &mirostat_mu);
    }
}

int llama_api_sample_token_greedy(void* container, void* candidates_a_p) {
    int id = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        id = llama_sample_token_greedy((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p);
    }

    return id;
}

int llama_api_sample_token(void* container, void* candidates_a_p) {
    int id = 0;

    variables_container* c = (variables_container*)container;
    if ((llama_context*)c->ctx != NULL) {
        id = (int)llama_sample_token((llama_context*)c->ctx, (llama_token_data_array*)candidates_a_p);
    }

    return id;
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
