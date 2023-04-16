package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"

func (l *LLama) SetPrompt(prompt string) {
	container := l.Container
	C.llama_set_params_prompt(container, C.CString(prompt))
}

func (l *LLama) SetUserInput(input string) {
	container := l.Container
	C.llama_set_user_input(container, C.CString(input))
}
