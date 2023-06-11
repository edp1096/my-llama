#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include "myllama.h"

/* Initialize before main loop */
void* bd_init_container();
bool bd_load_model(void* container);
// void bd_init_params(void* container);
bool bd_allocate_variables(void* container);

/* For main loop */
bool bd_predict_tokens(void* container);
bool bd_receive_input(void* container);
bool bd_append_input(void* container);
bool bd_wait_or_continue(void* container);
// int bd_get_embed_id(void* container, int index);
// char* bd_get_embed_string(void* container, int id);

/* Others */
bool bd_check_prompt_or_continue(void* container);
void bd_dropback_user_input(void* container);

#ifdef __cplusplus
}
#endif
