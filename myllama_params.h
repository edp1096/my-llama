#ifdef __cplusplus
extern "C" {
#endif

void init_gpt_params(void* container);
void init_context_params(void* container);

/* Setters */
void set_gptparams_n_threads(void* container, int value);
void set_gptparams_use_mlock(void* container, bool value);
void set_gptparams_n_predict(void* container, int value);

#ifdef __cplusplus
}
#endif
