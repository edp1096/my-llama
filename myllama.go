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
	C.set_model_path(l.Container, C.CString(modelFNAME))

	result := bool(C.load_model(l.Container))
	// result := bool(C.bd_load_model(l.Container))
	if !result {
		return fmt.Errorf("failed to load the model")
	}

	return nil
}

func (l *LLama) PredictTokens() bool {
	return bool(C.bd_predict_tokens(l.Container))
}

func (l *LLama) Free() {
	C.llama_api_free(l.Container)
}

func (l *LLama) GetNumRemain() int {
	return int(C.get_n_remain(l.Container))
}

func (l *LLama) GetEmbedSize() int {
	return int(C.get_embd_size(l.Container))
}

func (l *LLama) GetEmbedString(idx int) string {
	id := C.get_embed_id(l.Container, C.int(idx))
	embedCSTR := C.get_embed_string(l.Container, id)
	embedSTR := C.GoString(embedCSTR)

	return embedSTR
}

func (l *LLama) SetIsInteracting(isInteracting bool) {
	C.set_is_interacting(l.Container, C.bool(isInteracting))
}

/* Save/Load State or session */
func (l *LLama) SaveState(fname string) {
	C.save_state(l.Container, C.CString(fname))
}

func (l *LLama) LoadState(fname string) {
	C.load_state(l.Container, C.CString(fname))
}

func (l *LLama) SaveSession(fname string) {
	C.save_session(l.Container, C.CString(fname))
}

func (l *LLama) LoadSession(fname string) {
	C.load_session(l.Container, C.CString(fname))
}

func (l *LLama) CheckPromptOrContinue() bool {
	return bool(C.bd_check_prompt_or_continue(l.Container))
}

func (l *LLama) DropBackUserInput() {
	C.bd_dropback_user_input(l.Container)
}

/* LLAMA_API */

func (l *LLama) LlamaApiContextDefaultParams() {
	C.llama_api_context_default_params(l.Container)
}

func (l *LLama) LlamaApiMmapSupported() bool {
	return bool(C.llama_api_mmap_supported())
}
func (l *LLama) LlamaApiMlockSupported() bool {
	return bool(C.llama_api_mlock_supported())
}

// Initialize the llama + ggml backend
// Call once at the start of the program
func (l *LLama) LlamaApiInitBackend() {
	C.llama_api_init_backend()
}

func (l *LLama) LlamaApiTimeUs() int64 {
	return int64(C.llama_api_time_us())
}

// Various functions for loading a ggml llama model.
// Allocate (almost) all memory needed for the model.
func (l *LLama) LlamaApiInitFromFile(fname string) {
	C.llama_api_init_from_file(l.Container, C.CString(fname))
}

// Frees all allocated memory
func (l *LLama) LlamaApiFree() {
	C.llama_api_free(l.Container)
}

func (l *LLama) LlamaApiModelQuantize(fnameINP string, fnameOUT string, ftypeINT int, threads int) {
	C.llama_api_model_quantize(l.Container, C.CString(fnameINP), C.CString(fnameOUT), C.int(ftypeINT), C.int(threads))
}

// Apply a LoRA adapter to a loaded model
// path_base_model is the path to a higher quality model to use as a base for
// the layers modified by the adapter. Can be NULL to use the current loaded model.
// The model needs to be reloaded before applying a new adapter, otherwise the adapter
// will be applied on top of the previous one
func (l *LLama) LlamaApiApplyLoraFromFile(pathLora string, pathBaseModel string, threads int) {
	C.llama_api_apply_lora_from_file(l.Container, C.CString(pathLora), C.CString(pathBaseModel), C.int(threads))
}

// Sets the current rng seed.
func (l *LLama) LlamaApiSetRandomNumberGenerationSeed(seed int) {
	C.llama_api_set_rng_seed(l.Container, C.int(seed))
}

/* Todo: save/load state data, session */

// Run the llama inference to obtain the logits and probabilities for the next token.
// tokens + n_tokens is the provided batch of new tokens to process
// n_past is the number of tokens to use from previous eval calls
// Returns true on success
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

// Convert the provided text into tokens.
// The tokens pointer must be large enough to hold the resulting tokens.
// Returns the number of tokens on success, no more than n_max_tokens
// Returns a negative number on failure - the number of tokens that would have been returned
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

// Token logits obtained from the last call to llama_eval()
// The logits for the last token are stored in the last row
// Can be mutated in order to change the probabilities of the next token
// c struct myllama_container.n_logits, myllama_container.n_vocab
func (l *LLama) LlamaApiGetLogits() {
	C.llama_api_get_logits(l.Container)
}

// Get the embeddings for the input
// shape: [n_embd] (1-dimensional)
func (l *LLama) LlamaApiGetEmbeddings(embeddingSize int) []float64 {
	embeddings := C.llama_api_get_embeddings(l.Container)
	defer C.free(unsafe.Pointer(embeddings))

	embeddingSlice := make([]float64, embeddingSize)
	for i := 0; i < embeddingSize; i++ {
		embeddingSlice[i] = float64(*(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(&embeddings)) + uintptr(i)*unsafe.Sizeof(C.float(0)))))
	}

	return embeddingSlice
}

// Token Id -> String. Uses the vocabulary in the provided context
func (l *LLama) LlamaApiTokenToStr(token int32) string {
	return C.GoString(C.llama_api_token_to_str(l.Container, C.int(token)))
}

// Special tokens
func (l *LLama) LlamaApiTokenBOS() int {
	return int(C.llama_api_token_bos())
}

// Special tokens
func (l *LLama) LlamaApiTokenEOS() int {
	return int(C.llama_api_token_eos())
}

// Special tokens
func (l *LLama) LlamaApiTokenNL() int {
	return int(C.llama_api_token_nl())
}

// @details Repetition penalty described in CTRL academic paper https://arxiv.org/abs/1909.05858, with negative logit fix.
//
//	func (l *LLama) LlamaApiSampleRepetitionPenalty() {
//		return C.llama_api_sample_repetition_penalty(l.Container)
//	}
//
// @details Frequency and presence penalties described in OpenAI API https://platform.openai.com/docs/api-reference/parameter-details.
//
//	func (l *LLama) LlamaApiSampleFrequencyAndPresencePenalties() {
//		return C.llama_api_sample_frequency_and_presence_penalties(l.Container)
//	}

// @details Sorts candidate tokens by their logits in descending order and calculate probabilities based on logits.
func (l *LLama) LlamaApiSampleSoftmax() {
	C.llama_api_sample_softmax(l.Container)
}

// @details Top-K sampling described in academic paper "The Curious Case of Neural Text Degeneration" https://arxiv.org/abs/1904.09751
func (l *LLama) LlamaApiSampleTopK(topK int) {
	C.llama_api_sample_top_k(l.Container, C.int(topK))
}

// @details Nucleus sampling described in academic paper "The Curious Case of Neural Text Degeneration" https://arxiv.org/abs/1904.09751
func (l *LLama) LlamaApiSampleTopP(topP float64) {
	C.llama_api_sample_top_p(l.Container, C.float(topP))
}

// @details Tail Free Sampling described in https://www.trentonbricken.com/Tail-Free-Sampling/.
func (l *LLama) LlamaApiSampleTailFree(tfsZ float64) {
	C.llama_api_sample_tail_free(l.Container, C.float(tfsZ))
}

// @details Locally Typical Sampling implementation described in the paper https://arxiv.org/abs/2202.00666.
func (l *LLama) LlamaApiSampleTypical(typicalP float64) {
	C.llama_api_sample_typical(l.Container, C.float(typicalP))
}

func (l *LLama) LlamaApiSampleTemperature(temperature float64) {
	C.llama_api_sample_temperature(l.Container, C.float(temperature))
}

/* Todo: Mirostat v1 */

// @details Mirostat 2.0 algorithm described in the paper https://arxiv.org/abs/2007.14966. Uses tokens instead of words.
// @param candidates A vector of `llama_token_data` containing the candidate tokens, their probabilities (p), and log-odds (logit) for the current position in the generated text.
// @param tau  The target cross-entropy (or surprise) value you want to achieve for the generated text. A higher value corresponds to more surprising or less predictable text, while a lower value corresponds to less surprising or more predictable text.
// @param eta The learning rate used to update `mu` based on the error between the target and observed surprisal of the sampled word. A larger learning rate will cause `mu` to be updated more quickly, while a smaller learning rate will result in slower updates.
// @param mu Maximum cross-entropy. This value is initialized to be twice the target cross-entropy (`2 * tau`) and is updated in the algorithm based on the error between the target and observed surprisal.
func (l *LLama) LlamaApiSampleTokenMirostatV2(mirostatTAU, mirostatETA, mirostatMU float64) int32 {
	return int32(C.llama_api_sample_token_mirostat_v2(l.Container, C.float(mirostatTAU), C.float(mirostatETA), C.float(mirostatMU)))
}

// @details Selects the token with the highest probability.
func (l *LLama) LlamaApiSampleTokenGreedy() int32 {
	return int32(C.llama_api_sample_token_greedy(l.Container))
}

// @details Randomly selects a token from the candidates based on their probabilities.
func (l *LLama) LlamaApiSampleToken() int32 {
	return int32(C.llama_api_sample_token(l.Container))
}

// Print performance information
func (l *LLama) LlamaApiPrintTimings() {
	C.llama_api_print_timings(l.Container)
}

func (l *LLama) LlamaApiResetTimings() {
	C.llama_api_reset_timings(l.Container)
}

// Return system information string for screen printing
func (l *LLama) LlamaApiPrintSystemInfo() string {
	info := C.llama_api_print_system_info()

	return C.GoString(info)
}

/* Misc. */

func (l *LLama) AllocateTokens() {
	C.allocate_tokens(l.Container)
}

func (l *LLama) PrepareCandidates(numVocab int) {
	C.prepare_candidates(l.Container, C.int(numVocab))
}
