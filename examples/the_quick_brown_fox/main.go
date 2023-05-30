package main // import "the_quick_brown_fox"
import "C"
import (
	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	prompt := "The quick brown fox"

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.InitGptParams()
	l.InitContextParams()
	l.SetThreadsCount(4)
	l.SetUseMlock(true)

	// err = l.LoadModel(modelName)
	// if err != nil {
	// 	panic(err)
	// }

	l.LlamaApiInitFromFile(modelName)

	tokenCount := l.LlamaApiTokenize(prompt, true)

	println(tokenCount)

	// l.LlamaApiEval(text string, addBOS bool)
}
