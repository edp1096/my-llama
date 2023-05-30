#include "common.h"
#include "myllama.h"
#include "myllama_params.h"

void init_params(void* container) {
    variables_container* c = (variables_container*)container;
    gpt_params* params = (gpt_params*)c->params;

    llama_context_params* lparams = new llama_context_params(llama_context_default_params());

    lparams->n_ctx = params->n_ctx;
    lparams->n_gpu_layers = params->n_gpu_layers;
    lparams->seed = params->seed;
    lparams->f16_kv = params->memory_f16;
    lparams->use_mmap = params->use_mmap;
    lparams->use_mlock = params->use_mlock;
    lparams->logits_all = params->perplexity;
    lparams->embedding = params->embedding;

    c->params = lparams;
}
