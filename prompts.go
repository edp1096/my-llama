package myllama

// #include "binding.h"
// #include "myllama_params.h"
import "C"
import (
	"errors"
)

func (l *LLama) SetUserInput(input string) {
	container := l.Container
	C.set_user_input(container, C.CString(input))
}

func (l *LLama) AppendInput() (err error) {
	container := l.Container
	ok := bool(C.bd_append_input(container))
	if !ok {
		err = errors.New("failed to append the input")
	}

	return err
}

func (l *LLama) SetPrompt(prompt string) {
	container := l.Container
	C.set_gptparams_prompt(container, C.CString(prompt))
}

func (l *LLama) SetAntiPrompt(antiprompt string) {
	container := l.Container
	C.set_gptparams_antiprompt(container, C.CString(antiprompt))
}
