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
    int id;

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
    void* embd_inp;  // will be removed

    void* last_n_tokens;
    void* llama_token_newline;

    char* user_input;  // will be removed
};

void* init_container();

void allocate_tokens(void* container);
int* get_tokens(void* container);
void prepare_candidates(void* container, int n_vocab);
// int prepare_candidates(void* container, int n_vocab);

#ifdef __cplusplus
}
#endif
