package myllama

// #include "binding.h"
import "C"
import (
	"errors"
)

func (l *LLama) SetPrompt(prompt string) {
	container := l.Container
	C.bd_set_params_prompt(container, C.CString(prompt))
}

func (l *LLama) SetAntiPrompt(antiprompt string) {
	container := l.Container
	C.bd_set_params_antiprompt(container, C.CString(antiprompt))
}

func (l *LLama) SetUserInput(input string) {
	container := l.Container
	C.bd_set_user_input(container, C.CString(input))
}

func (l *LLama) AppendInput() (err error) {
	container := l.Container
	ok := bool(C.bd_append_input(container))
	if !ok {
		err = errors.New("failed to append the input")
	}

	return err
}
