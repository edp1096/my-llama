package main // import "run-cgollama"

import (
	"fmt"

	"myllama"
)

func main() {
	fmt.Println("Hello, My llama!")

	var err error

	l, err := myllama.New()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize the container: %s", err))
	}

	err = l.LoadModel("model.bin")
	if err != nil {
		panic(fmt.Sprintf("failed to load the model: %s", err))
	}
}
