package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"
import (
	"errors"
)

func (l *LLama) MakeReadyToPredict() (err error) {
	container := l.Container
	result := bool(C.bd_make_ready_to_predict(container))

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
func (l *LLama) SetThreadsCount(threads int) {
	C.bd_set_params_n_threads(l.Container, C.int(threads))
}

/* Setters - gpt_params / sampling parameters */
func (l *LLama) SetTopK(topK int) {
	C.bd_set_params_top_k(l.Container, C.int(topK))
}

func (l *LLama) InitParams() {
	C.bd_init_params(l.Container)
}
