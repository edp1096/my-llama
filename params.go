package myllama

// #include "binding.h"
// #include "myllama_params.h"
import "C"
import (
	"errors"
)

func (l *LLama) InitGptParams() {
	C.init_gpt_params(l.Container)
}
func (l *LLama) InitContextParams() {
	C.init_context_params(l.Container)
}

func (l *LLama) AllocateVariables() (err error) {
	container := l.Container
	result := bool(C.bd_allocate_variables(container))

	if !result {
		err = errors.New("failed to initialize the parameters")
	}

	return err
}

/* Getters - gpt_params */
func (l *LLama) GetThreadsCount() int {
	return int(C.bd_get_params_n_threads(l.Container))
}

/* Getters - gpt_params / sampling parameters */
func (l *LLama) GetTopK() int {
	return int(C.bd_get_params_top_k(l.Container))
}

func (l *LLama) GetTopP() float64 {
	return float64(C.bd_get_params_top_p(l.Container))
}

/* Setters - gpt_params */
func (l *LLama) SetNumThreads(threads int) {
	C.set_gptparams_n_threads(l.Container, C.int(threads))
}

func (l *LLama) SetUseMlock(useMlock bool) {
	C.set_gptparams_use_mlock(l.Container, C.bool(useMlock))
}

func (l *LLama) SetNumPredict(predicts int) {
	C.set_gptparams_n_predict(l.Container, C.int(predicts))
}

/* Setters - gpt_params / sampling parameters */
func (l *LLama) SetNCtx(nCtx int) {
	C.bd_set_params_n_ctx(l.Container, C.int(nCtx))
}
func (l *LLama) SetNBatch(nBatch int) {
	C.bd_set_params_n_batch(l.Container, C.int(nBatch))
}
func (l *LLama) SetSamplingMethod(method string) {
	mirostat := 0

	switch method {
	case "mirostat1":
		mirostat = 1
	case "mirostat2":
		mirostat = 2
	}

	C.bd_set_params_top_k(l.Container, C.int(mirostat))
}
func (l *LLama) SetTopK(topK int) {
	C.bd_set_params_top_k(l.Container, C.int(topK))
}
func (l *LLama) SetTopP(topP float64) {
	C.bd_set_params_top_p(l.Container, C.float(topP))
}
func (l *LLama) SetTemperature(temper float64) {
	C.bd_set_params_temper(l.Container, C.float(temper))
}
func (l *LLama) SetRepeatPenalty(penalty float64) {
	C.bd_set_params_repeat_penalty(l.Container, C.float(penalty))
}
