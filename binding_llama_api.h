#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

// void llama_api_context_default_params();

bool llama_api_mmap_supported();
bool llama_api_mlock_supported();
void llama_api_init_backend();

int64_t llama_api_time_us();
// void llama_api_init_from_file();
void llama_api_free(void* container);
// int llama_api_model_quantize();
// int llama_api_get_kv_cache_token_count();
void llama_api_set_rng_seed(void* container, int seed);
// int llama_api_get_state_size();
// int llama_api_copy_state_data();
// int llama_api_set_state_data();
// bool llama_api_load_session_file();
// bool llama_api_save_session_file();
// int llama_api_eval();
// int llama_api_tokenize();
int llama_api_n_vocab(void* container);
int llama_api_n_ctx(void* container);
int llama_api_n_embd(void* container);
// float* llama_api_get_logits(void* container);
// float* llama_api_get_embeddings();
// char* llama_api_token_to_str();
int llama_api_token_bos();
int llama_api_token_eos();
int llama_api_token_nl();
// void llama_api_sample_repetition_penalty();
// void llama_api_sample_frequency_and_presence_penalties();
// void llama_api_sample_softmax();
// void llama_api_sample_top_k();
// void llama_api_sample_top_p();
// void llama_api_sample_tail_free();
// void llama_api_sample_typical();
// void llama_api_sample_temperature();
// void llama_api_sample_token_mirostat_v2();
// void llama_api_sample_token_greedy();
// void llama_api_sample_token();
void llama_api_print_timings(void* container);
void llama_api_reset_timings(void* container);

char* llama_api_print_system_info();

#ifdef __cplusplus
}
#endif
