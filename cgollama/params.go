package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"
import (
	"errors"
)

func (l *LLama) InitParams() (err error) {
	container := l.Container
	result := bool(C.llama_init_params(container))

	if !result {
		err = errors.New("failed to initialize the parameters")
	}

	return err
}

func (l *LLama) SetupParams() {
	container := l.Container
	C.llama_setup_params(container)
}
