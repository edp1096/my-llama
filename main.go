package main // import "my-llama"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	llama "my-llama/cgollama"
	"os"
	"runtime"
	"strconv"
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

func initModel(model string) (*llama.LLama, llama.PredictOptions) {
	l, err := llama.New(model, llama.SetContext(128), llama.SetParts(-1))
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		os.Exit(1)
	}

	po := llama.NewPredictOptions(llama.SetTokens(tokens), llama.SetThreads(threads), llama.SetTopK(90), llama.SetTopP(0.86))
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}

	return l, po
}

func changePredictOptions(po *llama.PredictOptions, newTOKEN int, newThreads int) {
	po.Tokens = newTOKEN
	po.Threads = newThreads
}

func main() {
	var model string

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&model, "m", "./ggml-llama_7b-q4_1.bin", "path to quantized ggml model file to load")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	flags.IntVar(&tokens, "n", 256, "number of tokens to predict")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	l, po := initModel(model)

	fmt.Printf("Model loaded.\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		text := readMultiLineInput(reader)

		if strings.HasPrefix(text, "/") {
			text = strings.TrimSpace(text)
			text = strings.TrimPrefix(text, "/")
			splitText := strings.Split(text, "=")

			needChange := false
			newTOKEN := tokens
			newThreads := threads

			switch splitText[0] {
			case "tokens":
				newTOKEN, err = strconv.Atoi(splitText[1])
				if err != nil {
					fmt.Println("Error: ", err)
				}

				fmt.Println("Tokens: ", newTOKEN)
				needChange = true
			case "threads":
				newThreads, err = strconv.Atoi(splitText[1])
				if err != nil {
					fmt.Println("Error: ", err)
				}

				fmt.Println("Threads: ", newThreads)
				needChange = true
			default:
			}

			if needChange {
				changePredictOptions(&po, newTOKEN, newThreads)
			}
			continue
		}

		err := l.Predict(text, po)
		if err != nil {
			panic(err)
		}
	}
}
