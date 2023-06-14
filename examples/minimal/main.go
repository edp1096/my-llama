package main // import "minimal"

import (
	"fmt"

	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	numPredict := 16

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.LlamaApiInitBackend()
	l.InitGptParams()

	l.SetNumThreads(4)
	l.SetUseMlock(true)
	l.SetNumPredict(numPredict)
	l.SetNumGpuLayers(32)
	l.SetSeed(42)

	l.InitContextParamsFromGptParams()

	// l.LlamaApiInitFromFile(modelName)
	err = l.LoadModel(modelName)
	if err != nil {
		panic(err)
	}

	l.AllocateTokens()

	numPast := 0
	prompt := "The quick brown fox"

	promptTokens, promptNumTokens := l.LlamaApiTokenize(prompt, true)
	fmt.Println("promptTokens:", promptTokens)

	if promptNumTokens < 1 {
		fmt.Println("numToken < 1")
		panic("numToken < 1")
	}

	isOK := l.LlamaApiEval(promptTokens, promptNumTokens, numPast)
	numPast += promptNumTokens

	fmt.Println("n_prompt_token, n_past, isOK:", promptNumTokens, numPast, isOK)
	fmt.Println("numPredict:", numPredict)

	for i := 0; i < numPredict; i++ {
		l.LlamaApiGetLogits()
		numVocab := l.LlamaApiNumVocab()

		l.PrepareCandidates(numVocab)
		nextToken := l.LlamaApiSampleToken()
		nextTokenStr := l.LlamaApiTokenToStr(nextToken)

		fmt.Print(nextTokenStr)
		l.LlamaApiEval([]int32{nextToken}, 1, numPast)

		numPast++
	}

	fmt.Println()

	l.LlamaApiFree()
}
