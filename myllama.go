package myllama

/*
#cgo CXXFLAGS: -Illama.cpp -Illama.cpp/examples
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

	return nil
}

func (l *LLama) PredictTokens() bool {
	return bool(C.bd_predict_tokens(l.Container))
}

func (l *LLama) FreeParams() {
	C.bd_free_params(l.Container)
}

func (l *LLama) FreeModel() {
	C.bd_free_model(l.Container)
}

func (l *LLama) FreeALL() {
	l.FreeParams()
	l.FreeModel()
}

func (l *LLama) GetRemainCount() int {
	return int(C.bd_get_n_remain(l.Container))
}

func (l *LLama) GetEmbedSize() int {
	return int(C.bd_get_embd_size(l.Container))
}

func (l *LLama) GetEmbedString(idx int) string {
	id := C.bd_get_embed_id(l.Container, C.int(idx))
	embedCSTR := C.bd_get_embed_string(l.Container, id)
	embedSTR := C.GoString(embedCSTR)

	return embedSTR
}

func (l *LLama) SetIsInteracting(isInteracting bool) {
	C.bd_set_is_interacting(l.Container, C.bool(isInteracting))
}

func (l *LLama) SaveState(fname string) {
	C.bd_save_state(l.Container, C.CString(fname))
}

func (l *LLama) LoadState(fname string) {
	C.bd_load_state(l.Container, C.CString(fname))
}

func (l *LLama) SaveSession(fname string) {
	C.bd_save_session(l.Container, C.CString(fname))
}

func (l *LLama) LoadSession(fname string) {
	C.bd_load_session(l.Container, C.CString(fname))
}

func (l *LLama) CheckPromptOrContinue() bool {
	return bool(C.bd_check_prompt_or_continue(l.Container))
}

func (l *LLama) DropBackUserInput() {
	C.bd_dropback_user_input(l.Container)
}

func (l *LLama) PrintTimings() {
	C.bd_print_timings(l.Container)
}
