#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include "myllama.h"

/* Initialize before main loop */
void* bd_init_container();
bool bd_load_model(void* container);
void bd_init_params(void* container);
bool bd_allocate_variables(void* container);

/* For main loop */
bool bd_predict_tokens(void* container);
bool bd_receive_input(void* container);
bool bd_append_input(void* container);
bool bd_wait_or_continue(void* container);
int bd_get_embed_id(void* container, int index);
char* bd_get_embed_string(void* container, int id);

/* Finish loop */
void bd_free_params(void* container);
void bd_free_model(void* container);

/* Frees */
void bd_free_params(void* container);
void bd_free_model(void* container);

/* Getters */
int bd_get_n_remain(void* container);
int bd_get_params_n_predict(void* container);
bool bd_get_noecho(void* container);
int bd_get_embd_size(void* container);
int bd_get_embd_inp_size(void* container);
int bd_get_n_consumed(void* container);
bool bd_get_params_interactive_first(void* container);
bool bd_get_params_interactive(void* container);

/* Getters - gpt_params */
int bd_get_params_n_threads(void* container);

/* Getters - gpt_params / sampling parameters */
int bd_get_params_top_k(void* container);
float bd_get_params_top_p(void* container);
float bd_get_params_temper(void* container);
float bd_get_params_repeat_penalty(void* container);

/* Setters */
void bd_set_params_interactive_first(void* container);
void bd_set_is_interacting(void* container, bool is_interacting);
void bd_set_n_remain(void* container, int n_predict);
void bd_set_model_path(void* container, char* path);
void bd_set_params_antiprompt(void* container, char* antiprompt);
void bd_set_params_prompt(void* container, char* prompt);
void bd_set_user_input(void* container, const char* user_input);

/* Setters - gpt_params */
void bd_set_params_n_threads(void* container, int value);
void bd_set_params_use_mlock(void* container, bool value);

/* Setters - gpt_params / sampling parameters */
void bd_set_params_n_ctx(void* container, int value);
void bd_set_params_n_batch(void* container, int value);
void bd_set_sampling_method(void* container, int value);
void bd_set_params_top_k(void* container, int value);
void bd_set_params_top_p(void* container, float value);
void bd_set_params_temper(void* container, float value);
void bd_set_params_repeat_penalty(void* container, float value);

/* State dump */
void bd_save_state(void* container, char* fname);
void bd_load_state(void* container, char* fname);
void bd_save_session(void* container, char* fname);
void bd_load_session(void* container, char* fname);

/* Others */
bool bd_check_prompt_or_continue(void* container);
void bd_dropback_user_input(void* container);
void bd_print_timings(void* container);

#ifdef __cplusplus
}
#endif
