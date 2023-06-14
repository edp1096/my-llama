#ifdef __cplusplus
extern "C" {
#endif

void init_gpt_params(void* container);
void init_context_params_from_gpt_params(void* container);

/* Getters - gptparams */
int get_gptparams_n_threads(void* container);
int get_gptparams_top_k(void* container);
float get_gptparams_top_p(void* container);

/* Setters - gptparams */
void set_gptparams_seed(void* container, int value);
void set_gptparams_n_threads(void* container, int value);
void set_gptparams_use_mlock(void* container, bool value);
void set_gptparams_n_predict(void* container, int value);
void set_gptparams_prompt(void* container, char* prompt);
void set_gptparams_antiprompt(void* container, char* antiprompt);
void set_gptparams_n_gpu_layers(void* container, int value);
void set_gptparams_embedding(void* container, bool value);

/* Setters - gptparams / sampling parameters */
void set_gptparams_n_ctx(void* container, int value);
void set_gptparams_n_batch(void* container, int value);
// void set_gptparams_sampling_method(void* container, int value);
void set_gptparams_top_k(void* container, int value);
void set_gptparams_top_p(void* container, float value);
void set_gptparams_temperature(void* container, float value);
void set_gptparams_repeat_penalty(void* container, float value);

#ifdef __cplusplus
}
#endif
