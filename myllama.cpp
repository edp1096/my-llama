#include <cstring>
#include <algorithm>

#include "common.h"
#include "myllama.h"
#include "myllama_params.h"

void* init_container() {
    myllama_container* c = new myllama_container;

    c->gptparams = (void*)new gpt_params;
    c->ctxparams = (void*)new llama_context_params(llama_context_default_params());
    c->session_tokens = (void*)new std::vector<llama_token>;

    init_gpt_params(c);

    return c;
}

bool load_model(void* container) {
    myllama_container* c = (myllama_container*)container;
    gpt_params* gptparams = (gpt_params*)c->gptparams;
    llama_context_params* ctxparams = (llama_context_params*)c->ctxparams;

    if (gptparams->seed < 0) {
        gptparams->seed = time(NULL);
    }

    // gptparams->n_gpu_layers = 32;
    // printf("n_gpu_layers: %d\n", gptparams->n_gpu_layers);

    llama_init_backend();

    // init_context_params_from_gpt_params(c);

    printf("Model: %s\n", gptparams->model.c_str());
    llama_context* ctx = llama_init_from_file(gptparams->model.c_str(), *ctxparams);

    if (ctx == NULL) {
        fprintf(stderr, "%s: error: failed to load model '%s'\n", __func__, gptparams->model.c_str());
        return false;
    }

    if (!gptparams->lora_adapter.empty()) {
        int err = llama_apply_lora_from_file(
            ctx,
            gptparams->lora_adapter.c_str(),
            gptparams->lora_base.empty() ? NULL : gptparams->lora_base.c_str(),
            gptparams->n_threads);
        if (err != 0) {
            fprintf(stderr, "%s: error: failed to apply lora adapter\n", __func__);
            return false;
        }
    }

    c->ctx = ctx;

    return true;
}

void set_model_path(void* container, char* path) {
    myllama_container* c = (myllama_container*)container;

    ((gpt_params*)c->gptparams)->model = path;
}

void allocate_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;

    std::vector<llama_token> tokens(((gpt_params*)c->gptparams)->n_ctx);

    c->tokens = (void*)tokens.data();
    c->n_tokens = tokens.size();
}

int* get_tokens(void* container) {
    myllama_container* c = (myllama_container*)container;

    return (int*)c->tokens;
}

void prepare_candidates(void* container, int n_vocab) {
    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    float* logits = c->logits;

    llama_token_data_array* candidates_p = (llama_token_data_array*)malloc(sizeof(llama_token_data_array));
    candidates_p->data = (llama_token_data*)malloc(sizeof(llama_token_data) * n_vocab);
    candidates_p->size = n_vocab;
    candidates_p->sorted = false;

    for (llama_token token_id = 0; token_id < n_vocab; token_id++) {
        candidates_p->data[token_id] = llama_token_data{token_id, logits[token_id], 0.0f};
    }

    c->candidates = (void*)candidates_p;
}

/* Getters */
int get_n_remain(void* container) {
    return ((myllama_container*)container)->n_remain;
}

int get_embd_size(void* container) {
    myllama_container* c = (myllama_container*)container;

    return (int)((std::vector<llama_token>*)c->embd)->size();
}

int get_embed_id(void* container, int index) {
    myllama_container* c = (myllama_container*)container;

    return ((std::vector<llama_token>*)c->embd)->at(index);
}

char* get_embed_string(void* container, int id) {
    myllama_container* c = (myllama_container*)container;

    return const_cast<char*>(llama_token_to_str((llama_context*)c->ctx, id));
}

/* Setters */
void set_is_interacting(void* container, bool is_interacting) {
    myllama_container* c = (myllama_container*)container;

    c->is_interacting = is_interacting;
}

void set_user_input(void* container, const char* user_input) {
    myllama_container* c = (myllama_container*)container;

    c->user_input = strdup(user_input);
}

/* Misc. */
void save_state(void* container, char* fname) {
    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    size_t ctx_size = llama_get_state_size(ctx);
    uint8_t* state_mem = new uint8_t[ctx_size];
    llama_copy_state_data(ctx, state_mem);

    FILE* fp_write = fopen(fname, "wb");
    fwrite(&ctx_size, 1, sizeof(size_t), fp_write);
    fwrite(state_mem, 1, ctx_size, fp_write);
    fwrite(&c->last_n_tokens, 1, sizeof(int), fp_write);
    fwrite(&c->n_past, 1, sizeof(int), fp_write);
    fclose(fp_write);

    delete[] state_mem;
}

void load_state(void* container, char* fname) {
    myllama_container* c = (myllama_container*)container;

    FILE* fp_read = fopen(fname, "rb");
    size_t ctx_size = llama_get_state_size((llama_context*)c->ctx);

    size_t state_size;
    uint8_t* state_mem = new uint8_t[ctx_size];
    fread(&state_size, 1, sizeof(size_t), fp_read);

    if (state_size != ctx_size) {
        printf("Error: state size mismatch. Expected %zu, got %zu\n", ctx_size, state_size);
    }

    fread(state_mem, 1, state_size, fp_read);
    fread(&c->last_n_tokens, 1, sizeof(int), fp_read);
    fread(&c->n_past, 1, sizeof(int), fp_read);
    fclose(fp_read);

    llama_set_state_data((llama_context*)c->ctx, state_mem);
    delete[] state_mem;
}

void save_session(void* container, char* fname) {
    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    llama_save_session_file(
        ctx,
        fname,
        ((std::vector<llama_token>*)c->session_tokens)->data(),
        ((std::vector<llama_token>*)c->session_tokens)->size());
}

void load_session(void* container, char* fname) {
    myllama_container* c = (myllama_container*)container;
    llama_context* ctx = (llama_context*)c->ctx;

    // fopen to check for existing session
    FILE* fp = std::fopen(fname, "rb");
    if (fp != NULL) {
        std::fclose(fp);

        ((std::vector<llama_token>*)c->session_tokens)->resize(((gpt_params*)c->gptparams)->n_ctx);
        size_t n_token_count_out = 0;
        if (!llama_load_session_file(
                ctx,
                fname,
                ((std::vector<llama_token>*)c->session_tokens)->data(),
                ((std::vector<llama_token>*)c->session_tokens)->capacity(),
                &n_token_count_out)) {
            fprintf(stderr, "%s: error: failed to load session file '%s'\n", __func__, fname);
            return;
        }
        ((std::vector<llama_token>*)c->session_tokens)->resize(n_token_count_out);

        fprintf(stderr, "%s: loaded a session with prompt size of %d tokens\n", __func__, (int)((std::vector<llama_token>*)c->session_tokens)->size());
    } else {
        fprintf(stderr, "%s: session file does not exist, will create\n", __func__);
    }
}

/* From example/main.cpp */
