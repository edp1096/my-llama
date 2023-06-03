package main // import "the_quick_brown_fox"

import (
	"fmt"

	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	prompt := "The quick brown fox"

	// threadsCount := 4
	predictCount := 16

	numPast := 0

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.InitGptParams()
	l.InitContextParams()
	// l.SetThreadsCount(threadsCount)
	// l.SetUseMlock(true)
	// l.SetPredictCount(predictCount)

	l.LlamaApiInitFromFile(modelName)

	l.AllocateTokens()

	promptTokens, promptNumTokens := l.LlamaApiTokenize(prompt, true)
	fmt.Println("promptTokens:", promptTokens)

	if promptNumTokens < 1 {
		fmt.Println("tokenCount < 1")
		panic("tokenCount < 1")
	}

	isOK := l.LlamaApiEval(promptTokens, promptNumTokens, numPast)
	numPast += promptNumTokens

	fmt.Println("n_prompt_token, n_past, isOK:", promptNumTokens, numPast, isOK)
	fmt.Println("predictCount:", predictCount)

	for i := 0; i < predictCount; i++ {
		l.LlamaApiGetLogits()
		numVocab := 0
		numVocab = l.LlamaApiNumVocab()

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
