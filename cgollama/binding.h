#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

struct pred_variables {
    int n_past;
    int n_remain;
    int n_consumed;

    int n_ctx;

    void* embd;
    void* embedding_inp;
    void* context;

    void* last_n_tokens;
};

void* load_model(const char *fname, int n_ctx, int n_parts, int n_seed, bool memory_f16, bool mlock);

void* llama_allocate_params(const char* prompt, const char* antiprompt, int seed, int threads, int tokens,
                            int top_k, float top_p, float temp, float repeat_penalty, int repeat_last_n, bool ignore_eos, bool memory_f16);

void llama_free_params(void* params_ptr);

void llama_free_model(void* state);

void* llama_update_params(void* params_ptr, const char* prompt);
int llama_get_remain_count(void* pred_vars_v);
void llama_default_signal_action();
void* llama_prepare_pred_vars(void* params_ptr, void* state_pr);
int llama_get_embedding_ids(void* params_ptr, void* pred_vars_ptr);
int llama_get_id(void* pred_vars_ptr, int index);
char* llama_get_embed_string(void* pred_vars_ptr, int id);
bool llama_check_token_end(void* pred_vars_ptr);

int llama_predict(void* params_ptr, void* state_pr, char* result);

#ifdef __cplusplus
}
#endif
