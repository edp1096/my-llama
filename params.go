package myllama

// #include "myllama.h"
// #include "myllama_params.h"
import "C"
import (
	"errors"
)

func (l *LLama) InitGptParams() {
	C.init_gpt_params(l.Container)
}
func (l *LLama) InitContextParamsFromGptParams() {
	C.init_context_params_from_gpt_params(l.Container)
}

func (l *LLama) AllocateVariables() (err error) {
	container := l.Container
	result := bool(C.allocate_variables(container))

	if !result {
		err = errors.New("failed to initialize the parameters")
	}

	return err
}

/* Getters - gpt_params */
func (l *LLama) GetThreadsCount() int {
	return int(C.get_gptparams_n_threads(l.Container))
}

/* Getters - gpt_params / sampling parameters */
func (l *LLama) GetTopK() int {
	return int(C.get_gptparams_top_k(l.Container))
}

func (l *LLama) GetTopP() float64 {
	return float64(C.get_gptparams_top_p(l.Container))
}

/* Setters - gpt_params */
func (l *LLama) SetSeed(value int) {
	C.set_gptparams_seed(l.Container, C.int(value))
}

func (l *LLama) SetNumThreads(threads int) {
	C.set_gptparams_n_threads(l.Container, C.int(threads))
}

func (l *LLama) SetUseMlock(useMlock bool) {
	C.set_gptparams_use_mlock(l.Container, C.bool(useMlock))
}

func (l *LLama) SetNumPredict(predicts int) {
	C.set_gptparams_n_predict(l.Container, C.int(predicts))
}

func (l *LLama) SetNumGpuLayers(numLayers int) {
	C.set_gptparams_n_gpu_layers(l.Container, C.int(numLayers))
}

func (l *LLama) SetEmbedding(useEmbedding bool) {
	C.set_gptparams_embedding(l.Container, C.bool(useEmbedding))
}

/* Setters - gpt_params / sampling parameters */
func (l *LLama) SetNCtx(nCtx int) {
	C.set_gptparams_n_ctx(l.Container, C.int(nCtx))
}
func (l *LLama) SetNBatch(nBatch int) {
	C.set_gptparams_n_batch(l.Container, C.int(nBatch))
}
func (l *LLama) SetSamplingMethod(method string) {
	mirostat := 0

	switch method {
	case "mirostat1":
		mirostat = 1
	case "mirostat2":
		mirostat = 2
	}

	C.set_gptparams_top_k(l.Container, C.int(mirostat))
}
func (l *LLama) SetTopK(topK int) {
	C.set_gptparams_top_k(l.Container, C.int(topK))
}
func (l *LLama) SetTopP(topP float64) {
	C.set_gptparams_top_p(l.Container, C.float(topP))
}
func (l *LLama) SetTemperature(temper float64) {
	C.set_gptparams_temperature(l.Container, C.float(temper))
}
func (l *LLama) SetRepeatPenalty(penalty float64) {
	C.set_gptparams_repeat_penalty(l.Container, C.float(penalty))
}
