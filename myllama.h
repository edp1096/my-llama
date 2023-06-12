#ifdef __cplusplus
extern "C" {
#endif

struct tokens_array {
    int* data;
    int size;
    bool sorted;
};

struct myllama_container {
    void* ctx;

    void* gptparams;
    void* ctxparams;
    void* session_tokens;

    int n_past;
    int n_remain;
    int n_consumed;
    int n_session_consumed;

    bool is_interacting;
    bool input_noecho;
    int n_ctx;

    void* tokens;
    int n_tokens;

    float* logits;
    int n_logits;
    int n_vocab;

    void* candidates;

    void* embd;
    void* embd_inp;

    void* last_n_tokens;
    void* llama_token_newline;

    char* user_input;
};

void* init_container();

bool load_model(void* container);

void set_model_path(void* container, char* path);
void allocate_tokens(void* container);
int* get_tokens(void* container);
void prepare_candidates(void* container, int n_vocab);

/* Getters */
int get_n_remain(void* container);
int get_embd_size(void* container);
int get_embed_id(void* container, int index);
char* get_embed_string(void* container, int id);

/* Setters */
void set_is_interacting(void* container, bool is_interacting);
void set_user_input(void* container, const char* user_input);

/* Misc. */
void save_state(void* container, char* fname);
void load_state(void* container, char* fname);
void save_session(void* container, char* fname);
void load_session(void* container, char* fname);

#ifdef __cplusplus
}
#endif
