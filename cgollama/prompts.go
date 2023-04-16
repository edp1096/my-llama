package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"
import "errors"

func (l *LLama) SetPrompt(prompt string) {
	container := l.Container
	C.llama_set_params_prompt(container, C.CString(prompt))
}

func (l *LLama) SetAntiPrompt(antiprompt string) {
	container := l.Container
	C.llama_set_params_antiprompt(container, C.CString(antiprompt))
}

func (l *LLama) SetUserInput(input string) {
	container := l.Container
	C.llama_set_user_input(container, C.CString(input))
}

func (l *LLama) AppendInput() (err error) {
	container := l.Container
	ok := bool(C.llama_append_input(container))
	if !ok {
		err = errors.New("failed to append the input")
	}

	return err
}
