package main // import "embedding"

import "C"
import (
	"fmt"

	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.LlamaApiInitBackend()
	l.InitGptParams()

	l.SetNumThreads(4)
	l.SetUseMlock(true)
	l.SetNumGpuLayers(32)
	l.SetEmbedding(true)

	l.InitContextParamsFromGptParams()

	err = l.LoadModel(modelName)
	if err != nil {
		panic(err)
	}

	numPast := 0
	prompt := " " + "Hello world!"

	tokenData, tokenSize := l.LlamaApiTokenize(prompt, true)
	for i := 0; i < tokenSize; i++ {
		tokenString := l.LlamaApiTokenToStr(tokenData[i])
		fmt.Printf("Token %d -> %s\n", tokenData[i], tokenString)
	}

	isOK := l.LlamaApiEval(tokenData, tokenSize, numPast)
	if !isOK {
		panic("Eval failed")
	}

	numEmbedding := l.LlamaApiNumEmbd()
	fmt.Println("Embedding count:", numEmbedding)

	embeddings := l.LlamaApiGetEmbeddings(numEmbedding)
	fmt.Println("Embeddings:", embeddings)
}
