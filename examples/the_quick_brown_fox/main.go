package main // import "the_quick_brown_fox"
import "C"
import (
	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.InitParams()
	l.SetThreadsCount(4)
	l.SetUseMlock(true)

	// err = l.LoadModel(modelName)
	// if err != nil {
	// 	panic(err)
	// }

	l.LlamaApiInitFromFile(modelName)

	l.LlamaApiTokenize("Hello world!", true)
}
