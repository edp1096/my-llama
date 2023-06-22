package main

import (
	"errors"

	llama "github.com/edp1096/my-llama"
)

var lm *llama.LLama

type ModelExecutionConfig struct {
	ModelName    string
	NumThreads   int
	UseMlock     bool
	NumGpuLayers int
	Seed         int
}

// Todo: Add top, topk, sampling method...
type SamplingParameters struct {
	NumPredict int
}

func setupLLM(config ModelExecutionConfig, params SamplingParameters) error {
	var err error

	lm, err = llama.New()
	if err != nil {
		return err
	}

	lm.LlamaApiInitBackend()
	lm.InitGptParams()

	lm.SetNumThreads(config.NumThreads)
	lm.SetUseMlock(config.UseMlock)
	lm.SetNumPredict(params.NumPredict)
	lm.SetNumGpuLayers(config.NumGpuLayers)
	lm.SetSeed(config.Seed)

	lm.InitContextParamsFromGptParams()

	err = lm.LoadModel(config.ModelName)
	if err != nil {
		return err
	}

	lm.AllocateTokens()

	return err
}

func evalPrompt(prompt string, numPast int) (int, error) {
	var err error
	promptTokens, numPromptTokens := lm.LlamaApiTokenize(prompt, true)

	ok := lm.LlamaApiEval(promptTokens, numPromptTokens, numPast)
	numPast += numPromptTokens

	if !ok {
		err = errors.New("eval failed")
	}

	return numPast, err
}

func getTokenString(numPast int) (string, int) {
	lm.LlamaApiGetLogits()
	numVocab := lm.LlamaApiNumVocab()

	lm.PrepareCandidates(numVocab)

	nextToken := lm.LlamaApiSampleToken()
	nextTokenStr := lm.LlamaApiTokenToStr(nextToken)

	lm.LlamaApiEval([]int32{nextToken}, 1, numPast)
	numPast++

	return nextTokenStr, numPast
}
