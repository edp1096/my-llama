#ifdef __cplusplus
extern "C" {
#endif

struct myllama_container {
    int n_past;
    int n_remain;
    int n_consumed;
    int n_session_consumed;
    int id;

    bool is_interacting;
    bool input_noecho;
    int n_ctx;

    void* session_tokens;

    void* last_n_tokens;
    void* llama_token_newline;

    void* embd;
    void* embd_inp;

    void* ctx;
    void* gptparams;
    void* ctxparams;

    char* user_input;
};

#ifdef __cplusplus
}
#endif
