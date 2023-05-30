#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

void* llama_api_context_default_params();

bool llama_api_mmap_supported();
bool llama_api_mlock_supported();
void llama_api_init_backend();

int64_t llama_api_time_us();
void* llama_api_init_from_file(char* path_model, void* params_p);
void llama_api_free(void* container);

int llama_api_model_quantize(char* fname_inp, char* fname_out, int ftype, int nthread);
int llama_api_get_kv_cache_token_count(void* container);
void llama_api_set_rng_seed(void* container, int seed);

int llama_api_get_state_size(void* container);
int llama_api_copy_state_data(void* container, void* dst_p);
int llama_api_set_state_data(void* container, void* src_p);
bool llama_api_load_session_file(void* container, char* path_session, void* tokens_out, int n_token_capacity, int* n_token_count_out);
bool llama_api_save_session_file(void* container, char* path_session, void* tokens, int n_token_count);

int llama_api_eval(void* container, int* tokens, int n_tokens, int n_past, int n_threads);
int llama_api_tokenize(void* container, char* text, bool add_bos);
int llama_api_n_vocab(void* container);
int llama_api_n_ctx(void* container);
int llama_api_n_embd(void* container);

void* llama_api_get_logits(void* container);
float* llama_api_get_embeddings(void* container);
char* llama_api_token_to_str(void* container, int token);
int llama_api_token_bos();
int llama_api_token_eos();
int llama_api_token_nl();

void llama_api_sample_repetition_penalty(void* container, void* candidates_a_p, void* last_tokens, int last_tokens_size, float penalty);
void llama_api_sample_frequency_and_presence_penalties(void* container, void* candidates_a_p, void* last_tokens, int last_tokens_size, float alpha_frequency, float alpha_presence);
void llama_api_sample_softmax(void* container, void* candidates_a_p);
void llama_api_sample_top_k(void* container, void* candidates_a_p, int top_k);
void llama_api_sample_top_p(void* container, void* candidates_a_p, float top_p);
void llama_api_sample_tail_free(void* container, void* candidates_a_p, float tfs_z);
void llama_api_sample_typical(void* container, void* candidates_a_p, float typical_p);
void llama_api_sample_temperature(void* container, void* candidates_a_p, float temperature);
void llama_api_sample_token_mirostat_v2(void* container, void* candidates_a_p, float mirostat_tau, float mirostat_eta, float mirostat_mu);
int llama_api_sample_token_greedy(void* container, void* candidates_a_p);
int llama_api_sample_token(void* container, void* candidates_a_p);

void llama_api_print_timings(void* container);
void llama_api_reset_timings(void* container);

char* llama_api_print_system_info();

#ifdef __cplusplus
}
#endif
