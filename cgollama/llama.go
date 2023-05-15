package cgollama

// #cgo CFLAGS: -I../llama.cpp
// #cgo LDFLAGS: -static -L../ -lstdc++ -lllama
// #include "llama.h"
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
	return &LLama{}, nil
}

func (l *LLama) Hello() {
	fmt.Println("Hello, I'm LLama!")
}
