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
	result := bool(C.llama_make_ready_to_predict(container))

	if !result {
		err = errors.New("failed to initialize the parameters")
	}

	return err
}

func (l *LLama) SetThreadsCount(threads int) {
	container := l.Container
	C.llama_set_params_n_threads(container, C.int(threads))
}

func (l *LLama) SetTopK(topK int) {
	container := l.Container
	C.llama_set_params_top_k(container, C.int(topK))
}

func (l *LLama) InitParams() {
	container := l.Container
	C.llama_init_params(container)
}
