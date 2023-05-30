#ifdef __cplusplus
extern "C" {
#endif

void init_params(void* container);

/* Setters */
void set_params_n_threads(void* container, int value);
void set_params_use_mlock(void* container, bool value);

#ifdef __cplusplus
}
#endif
