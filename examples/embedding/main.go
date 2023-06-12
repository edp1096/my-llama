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

	l.InitContextParamsFromGptParams()
	l.SetNumThreads(4)
	l.SetUseMlock(true)

	err = l.LoadModel(modelName)
	if err != nil {
		panic(err)
	}

	l.LlamaApiTokenize("Hello world!", true)

	embdCount := l.LlamaApiNumEmbd()
	fmt.Println("Embedding count:", embdCount)

	embeddings := l.LlamaApiGetEmbeddings(embdCount)
	fmt.Println("Embeddings:", embeddings)
}
