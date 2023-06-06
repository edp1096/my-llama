#include "common.h"
#include "myllama.h"
#include "myllama_params.h"

void init_gpt_params(void* container) {
    myllama_container* c = (myllama_container*)container;
    gpt_params* gptparams = (gpt_params*)c->gptparams;

    // running setting
    gptparams->seed = -1;         // RNG seed
    gptparams->n_threads = 4;     // CPU threads
    gptparams->n_predict = -1;    // new tokens to predict
    gptparams->n_ctx = 512;       // context size
    gptparams->n_batch = 512;     // batch size for prompt processing (must be >=32 to use BLAS)
    gptparams->n_keep = 0;        // number of tokens to keep from initial prompt
    gptparams->n_gpu_layers = 0;  // number of layers to store in VRAM

    // sampling parameters
    gptparams->top_k = 40;                 // <= 0 to use vocab size
    gptparams->top_p = 0.95f;              // 1.0 = disabled
    gptparams->tfs_z = 1.00f;              // 1.0 = disabled
    gptparams->typical_p = 1.00f;          // 1.0 = disabled
    gptparams->temp = 0.80f;               // 1.0 = disabled
    gptparams->repeat_penalty = 1.10f;     // 1.0 = disabled
    gptparams->repeat_last_n = 64;         // last n tokens to penalize (0 = disable penalty, -1 = context size)
    gptparams->frequency_penalty = 0.00f;  // 0.0 = disabled
    gptparams->presence_penalty = 0.00f;   // 0.0 = disabled
    gptparams->mirostat = 0;               // 0 = disabled, 1 = mirostat, 2 = mirostat 2.0
    gptparams->mirostat_tau = 5.00f;       // target entropy
    gptparams->mirostat_eta = 0.10f;       // learning rate

    // gptparams->interactive = true;
    // gptparams->interactive_first = gptparams->interactive;
    // gptparams->antiprompt = {};

    // gptparams->n_predict = 512;
    // gptparams->use_mmap = false;
    // gptparams->use_mlock = true;

    gptparams->n_gpu_layers = 33;
    gptparams->seed = 42;
    gptparams->n_threads = 4;
    gptparams->n_predict = 16;
    // gptparams->repeat_last_n = 64;
    // gptparams->prompt = "The quick brown fox ";

    // gptparams->prompt.insert(0, 1, ' ');
}

void init_context_params(void* container) {
    myllama_container* c = (myllama_container*)container;
    gpt_params* gptparams = (gpt_params*)c->gptparams;
    llama_context_params* ctxparams = (llama_context_params*)c->ctxparams;

    ctxparams->n_ctx = gptparams->n_ctx;
    ctxparams->n_gpu_layers = gptparams->n_gpu_layers;
    ctxparams->seed = gptparams->seed;
    ctxparams->f16_kv = gptparams->memory_f16;
    ctxparams->use_mmap = gptparams->use_mmap;
    ctxparams->use_mlock = gptparams->use_mlock;
    // ctxparams->logits_all = gptparams->perplexity;
    // ctxparams->embedding = gptparams->embedding;
}

/* Setters - gptparams. require restart */
void set_gptparams_n_threads(void* container, int value) {
    ((gpt_params*)((myllama_container*)container)->gptparams)->n_threads = value;
}

void set_gptparams_use_mlock(void* container, bool value) {
    ((gpt_params*)((myllama_container*)container)->gptparams)->use_mlock = value;
}

void set_gptparams_n_predict(void* container, int value) {
    ((gpt_params*)((myllama_container*)container)->gptparams)->n_predict = value;
}
