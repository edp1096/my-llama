package myllama

// #include "myllama.h"
// #include "myllama_params.h"
import "C"
import (
	"errors"
)

func (l *LLama) SetUserInput(input string) {
	C.set_user_input(l.Container, C.CString(input))
}

func (l *LLama) AppendInput() (err error) {
	ok := bool(C.append_input(l.Container))
	if !ok {
		err = errors.New("failed to append the input")
	}

	return err
}

func (l *LLama) SetPrompt(prompt string) {
	C.set_gptparams_prompt(l.Container, C.CString(prompt))
}

func (l *LLama) SetAntiPrompt(antiprompt string) {
	C.set_gptparams_antiprompt(l.Container, C.CString(antiprompt))
}
