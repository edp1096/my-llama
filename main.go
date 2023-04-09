package main // import "my-llama"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	llama "my-llama/cgollama"
	"os"
	"runtime"
	"strings"
)

var (
	threads = 4
	tokens  = 128
)

// readMultiLineInput reads input until an empty line is entered.
func readMultiLineInput(reader *bufio.Reader) string {
	var lines []string
	fmt.Print("> ")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			fmt.Printf("Reading the prompt failed: %s", err)
			os.Exit(1)
		}

		if len(strings.TrimSpace(line)) == 0 {
			break
		}

		lines = append(lines, line)
	}

	text := strings.Join(lines, "")

	return text
}

func main() {
	var model string

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&model, "m", "./models/7B/ggml-model-q4_0.bin", "path to q4_0.bin model file to load")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	flags.IntVar(&tokens, "n", 512, "number of tokens to predict")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}
	l, err := llama.New(model, llama.SetContext(128), llama.SetParts(-1))
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		os.Exit(1)
	}

	po := llama.NewPredictOptions(llama.SetTokens(tokens), llama.SetThreads(threads), llama.SetTopK(90), llama.SetTopP(0.86))
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}

	fmt.Printf("Model loaded successfully.\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		text := readMultiLineInput(reader)
		// fmt.Printf("text: %s\n", text)

		_, err := l.Predict(text, po)
		if err != nil {
			panic(err)
		}

		// fmt.Println()

		// fmt.Printf("response: %s\n\n", res)

		// res, err := l.Predict(text, po)
		// if err != nil {
		// 	panic(err)
		// }
		// _, err := l.Predict(text, llama.SetTokens(tokens), llama.SetThreads(threads), llama.SetTopK(90), llama.SetTopP(0.86))
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("response: %s\n\n", res)
	}
}
