#include "common.h"
#include "myllama.h"
#include "myllama_params.h"

void init_gpt_params(void* container) {
    myllama_container* c = (myllama_container*)container;
    gpt_params* gptparams = (gpt_params*)c->gptparams;

    // gptparams->interactive = true;
    // gptparams->interactive_first = gptparams->interactive;
    // gptparams->antiprompt = {};

    // gptparams->n_predict = 512;
    // gptparams->use_mmap = false;
    // gptparams->use_mlock = true;

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
    // ctxparams->n_gpu_layers = gptparams->n_gpu_layers;
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
