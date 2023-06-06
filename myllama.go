package myllama

/*
#cgo CXXFLAGS: -Ivendors/llama.cpp -Ivendors/llama.cpp/examples
#include <stdlib.h>
#include "binding.h"
#include "myllama_llama_api.h"
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
	container := C.init_container()
	if container == nil {
		return nil, fmt.Errorf("failed to initialize the container")
	}

	sysInfo := C.llama_api_print_system_info()
	fmt.Printf("System Info: %s\n", C.GoString(sysInfo))

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

// SaveState - Not use.
func (l *LLama) SaveState(fname string) {
	C.bd_save_state(l.Container, C.CString(fname))
}

// LoadState - Not use.
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

/* LLAMA_API */

func (l *LLama) LlamaApiTimeUs() int64 {
	return int64(C.llama_api_time_us())
}

func (l *LLama) LlamaApiInitFromFile(fname string) {
	C.llama_api_init_from_file(l.Container, C.CString(fname))
}

func (l *LLama) LlamaApiFree() {
	C.llama_api_free(l.Container)
}

func (l *LLama) LlamaApiSetRandomNumberGenerationSeed(seed int) {
	C.llama_api_set_rng_seed(l.Container, C.int(seed))
}

func (l *LLama) LlamaApiEval(tokens []int32, tokenCount int, numPast int) (result bool) {
	result = false

	threadsCount := l.GetThreadsCount()
	tokensPtr := &tokens[0]

	isFail := int(C.llama_api_eval(l.Container, (*C.int)(unsafe.Pointer(tokensPtr)), C.int(tokenCount), C.int(numPast), C.int(threadsCount)))
	if isFail == 0 {
		result = true
	}

	return result
}

func (l *LLama) LlamaApiTokenize(text string, addBOS bool) ([]int32, int) {
	tokenSize := int(C.llama_api_tokenize(l.Container, C.CString(text), C.bool(addBOS)))
	tokenPtr := C.get_tokens(l.Container)
	defer C.free(unsafe.Pointer(tokenPtr))

	tokens := make([]int32, tokenSize)
	for i := 0; i < tokenSize; i++ {
		tokens[i] = int32(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(tokenPtr)) + uintptr(i)*unsafe.Sizeof(C.int(0)))))
	}

	return tokens, tokenSize
}

func (l *LLama) LlamaApiNumVocab() int {
	return int(C.llama_api_n_vocab(l.Container))
}

func (l *LLama) LlamaApiNumCtx() int {
	return int(C.llama_api_n_ctx(l.Container))
}

func (l *LLama) LlamaApiNumEmbd() int {
	return int(C.llama_api_n_embd(l.Container))
}

func (l *LLama) LlamaApiGetLogits() {
	C.llama_api_get_logits(l.Container)
}

func (l *LLama) LlamaApiGetEmbeddings(embeddingSize int) []float64 {
	embeddings := C.llama_api_get_embeddings(l.Container)
	defer C.free(unsafe.Pointer(embeddings))

	embeddingSlice := make([]float64, embeddingSize)
	for i := 0; i < embeddingSize; i++ {
		embeddingSlice[i] = float64(*(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(&embeddings)) + uintptr(i)*unsafe.Sizeof(C.float(0)))))
	}

	return embeddingSlice
}

func (l *LLama) LlamaApiTokenToStr(token int32) string {
	return C.GoString(C.llama_api_token_to_str(l.Container, C.int(token)))
}

func (l *LLama) LlamaApiTokenBOS() int {
	return int(C.llama_api_token_bos())
}

func (l *LLama) LlamaApiTokenEOS() int {
	return int(C.llama_api_token_eos())
}

func (l *LLama) LlamaApiTokenNL() int {
	return int(C.llama_api_token_nl())
}

//	func (l *LLama) LlamaApiSampleRepetitionPenalty() {
//		return C.llama_api_sample_repetition_penalty(l.Container)
//	}
//
//	func (l *LLama) LlamaApiSampleFrequencyAndPresencePenalties() {
//		return C.llama_api_sample_frequency_and_presence_penalties(l.Container)
//	}
func (l *LLama) LlamaApiSampleSoftmax() {
	C.llama_api_sample_softmax(l.Container)
}
func (l *LLama) LlamaApiSampleTopK(topK int) {
	C.llama_api_sample_top_k(l.Container, C.int(topK))
}
func (l *LLama) LlamaApiSampleTopP(topP float64) {
	C.llama_api_sample_top_p(l.Container, C.float(topP))
}
func (l *LLama) LlamaApiSampleTailFree(tfsZ float64) {
	C.llama_api_sample_tail_free(l.Container, C.float(tfsZ))
}
func (l *LLama) LlamaApiSampleTypical(typicalP float64) {
	C.llama_api_sample_typical(l.Container, C.float(typicalP))
}
func (l *LLama) LlamaApiSampleTemperature(temperature float64) {
	C.llama_api_sample_temperature(l.Container, C.float(temperature))
}

func (l *LLama) LlamaApiSampleTokenMirostatV2(mirostatTAU, mirostatETA, mirostatMU float64) int32 {
	return int32(C.llama_api_sample_token_mirostat_v2(l.Container, C.float(mirostatTAU), C.float(mirostatETA), C.float(mirostatMU)))
}
func (l *LLama) LlamaApiSampleTokenGreedy() int32 {
	return int32(C.llama_api_sample_token_greedy(l.Container))
}
func (l *LLama) LlamaApiSampleToken() int32 {
	return int32(C.llama_api_sample_token(l.Container))
}

/* Misc. */

func (l *LLama) AllocateTokens() {
	C.allocate_tokens(l.Container)
}

func (l *LLama) PrepareCandidates(numVocab int) {
	C.prepare_candidates(l.Container, C.int(numVocab))
}
