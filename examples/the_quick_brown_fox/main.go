package main // import "the_quick_brown_fox"
import "C"
import (
	"fmt"

	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	prompt := "The quick brown fox"

	threadsCount := 4
	predictCount := 16

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.InitGptParams()
	l.InitContextParams()
	l.SetThreadsCount(threadsCount)
	l.SetUseMlock(true)
	l.SetPredictCount(predictCount)

	l.LlamaApiInitFromFile(modelName)

	tokens, tokenCount := l.LlamaApiTokenize(prompt, true)

	fmt.Println(tokens)
	fmt.Println(tokenCount)

	l.LlamaApiEval()

	fmt.Println("predictCount:", predictCount)
	for i := 0; i < predictCount; i++ {
		l.LlamaApiGetLogits()
		numVocab := l.LlamaApiNumVocab()
		fmt.Println("numVocab:", numVocab)
	}

	l.LlamaApiFree()
}
