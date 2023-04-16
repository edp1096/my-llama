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
bool llama_load_model(void* container);
void llama_init_params(void* container);
bool llama_make_ready_to_predict(void* container);

/* For main loop */
bool llama_predict_tokens(void* container);
bool llama_receive_input(void* container);
bool llama_append_input(void* container);
bool llama_wait_or_continue(void* container);
int llama_get_embed_id(void* container, int index);
char* llama_get_embed_string(void* container, int id);

/* Finish loop */
void llama_free_params(void* container);
void llama_free_model(void* container);

/* Getters */
int llama_get_n_remain(void* container);
int llama_get_params_n_predict(void* container);
bool llama_get_noecho(void* container);
int llama_get_embd_size(void* container);
int llama_get_embd_inp_size(void* container);
int llama_get_n_consumed(void* container);
bool llama_get_params_interactive_start(void* container);
bool llama_get_params_interactive(void* container);

/* Setters */
void llama_set_params_interactive_start(void* container);
void llama_set_is_interacting(void* container, bool is_interacting);
void llama_set_params_n_remain(void* container, int n_predict);
void llama_set_model_path(void* container, char* path);
void llama_set_params_antiprompt(void* container, char* antiprompt);
void llama_set_params_prompt(void* container, char* prompt);
void llama_set_user_input(void* container, const char* user_input);

/* Others */
bool llama_check_prompt_or_continue(void* container);
void llama_dropback_user_input(void* container);

#ifdef __cplusplus
}
#endif
