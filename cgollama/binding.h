#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

struct variables_container {
    int n_past;
    int n_remain;
    int n_consumed;
    int id;

    bool is_interacting;
    bool input_noecho;
    int n_ctx;

    void* last_n_tokens;
    void* llama_token_newline;

    void* embd;
    void* embd_inp;

    void* ctx;
    void* params;

    char* user_input;
};

/* Initialize before main loop */
void* llama_init_container();
void llama_setup_params(void* container);
bool llama_initialize(void* container);

/* For main loop */
bool llama_predict_tokens(void* container);
char* llama_get_embed_string(void* container, int id);
bool llama_wait_or_continue(void* container);
int llama_get_n_remain(void* container);
int llama_get_embd_inp_size(void* container);
int llama_get_n_consumed(void* container);

#ifdef __cplusplus
}
#endif
