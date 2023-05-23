package myllama

/*
#cgo CFLAGS: -I./llama.cpp -I./llama.cpp/examples
#cgo CXXFLAGS: -I./llama.cpp -I./llama.cpp/examples
#cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
#include "llama.h"
#include "binding.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type LLama struct {
	Container      unsafe.Pointer
	PredictStop    chan bool
	Threads        int
	UseDumpSession bool
}

func New() (*LLama, error) {
	container := C.bd_init_container()
	if container == nil {
		return nil, fmt.Errorf("failed to initialize the container")
	}

	return &LLama{Container: container}, nil
}

func (l *LLama) LoadModel(modelFNAME string) error {
	C.bd_set_model_path(l.Container, C.CString(modelFNAME))

	result := bool(C.bd_load_model(l.Container))
	if !result {
		return fmt.Errorf("failed to load the model")
	}

	C.llama_print_system_info()

	return nil
}
