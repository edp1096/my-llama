package main // import "my-llama"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unsafe"

	_ "embed"

	ws "golang.org/x/net/websocket"

	llama "my-llama/cgollama"
)

//go:embed html/index.html
var html string

var (
	Address      string = "localhost"
	Port         string = "1323"
	HttpProtocol string = "http"
)

var (
	threads = 4
	tokens  = 256
)

var model string

func wsController(w http.ResponseWriter, req *http.Request) {
	ws.Handler(func(conn *ws.Conn) {
		defer conn.Close()

		requestCount := 0 // to ignore prompt

		l, po := initModel(model)
		fmt.Printf("Model loaded.\n")

		handler := ws.Message
		l.PredictStop = make(chan bool)

		disconnected := false
		predictRunning := false
		for !disconnected {
			// Read
			message := ""
			err := handler.Receive(conn, &message)
			if err != nil {
				fmt.Println("Receive error:", err)
				disconnected = true
			}

			// Print received message
			if len(message) > 0 {
				if message == "$$__STOP__$$" {
					if predictRunning {
						l.PredictStop <- true
					}
					continue
				}

				// fmt.Printf("%s\n", message) // Print received message from client
			}

			if predictRunning {
				continue
			}

			datas := strings.Split(message, "\n$$__SEPARATOR__$$\n")

			params := unsafe.Pointer(nil)
			predVARs := unsafe.Pointer(nil)
			remainCOUNT := 0

			// Predict and Write
			go func() {
				predictRunning = true

				if len(datas) < 3 {
					return
				}

				prompt := datas[1]
				input := prompt + datas[0]
				antiprompt := datas[2]

				fmt.Println("Input:", input)
				fmt.Println("Prompt:", prompt)
				fmt.Println("Antiprompt:", antiprompt)

				if requestCount == 0 {
					params, predVARs, remainCOUNT = l.GetInitialParams(input, prompt, antiprompt, po)
				} else {
					input = datas[0]
					params, predVARs, remainCOUNT = l.GetContinueParams(input, antiprompt, params, predVARs, po)
				}

				err = l.Predict(conn, handler, params, predVARs, remainCOUNT)
				if err != nil {
					fmt.Println("Predict error:" + err.Error())
					disconnected = true
				}

				predictRunning = false
				requestCount++
				fmt.Println("Request count:", requestCount)
			}()
		}
	}).ServeHTTP(w, req)
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

// func changePredictOptions(po *llama.PredictOptions, newTOKEN int, newThreads int) {
// 	po.Tokens = newTOKEN
// 	po.Threads = newThreads
// }

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&model, "m", "./ggml-llama_7b-q4_1.bin", "path to quantized ggml model file to load")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	flags.IntVar(&tokens, "n", 256, "number of tokens to predict")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(model); os.IsNotExist(err) {
		fmt.Printf("Model file %s does not exist", model)
		os.Exit(1)
	}

	listenURI := Address + ":" + Port
	uri := HttpProtocol + "://" + Address + ":" + Port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})
	http.HandleFunc("/ws", wsController)

	fmt.Printf("Server is running at %s\n", uri)

	switch os := runtime.GOOS; os {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "linux":
		exec.Command("xdg-open", uri).Start()
	case "darwin":
		exec.Command("open", uri).Start()
	default:
		fmt.Printf("%s: unsupported platform", os)
	}

	log.Fatal(http.ListenAndServe(listenURI, nil))
}
