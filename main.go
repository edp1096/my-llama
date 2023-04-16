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
	// tokens  = 256
)

var modelFNAME string

func wsController(w http.ResponseWriter, req *http.Request) {
	ws.Handler(func(conn *ws.Conn) {
		defer conn.Close()

		var err error

		requestCount := 0 // to ignore prompt

		l, err := llama.New()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = l.LoadModel(modelFNAME)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		l.InitParams()

		fmt.Println("Model initialized..")

		reflectionPrompt := ""
		input := ""
		antiprompt := ""

		handler := ws.Message
		l.PredictStop = make(chan bool)

		disconnected := false
		predictRunning := false
		for !disconnected {
			// Read
			message := ""
			err = handler.Receive(conn, &message)
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

			// Predict and Write
			go func() {
				predictRunning = true

				if len(datas) < 3 {
					return
				}

				reflectionPrompt = datas[1]
				antiprompt = datas[2]
				// input = antiprompt + datas[0]
				input = datas[0]

				if requestCount == 0 {
					l.SetPrompt(reflectionPrompt)
					l.SetAntiPrompt(antiprompt)

					err = l.MakeReadyToPredict()
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}

				l.SetUserInput(input)
				err = l.AppendInput()
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				l.SetIsInteracting(false)

				err = l.Predict(conn, handler)
				if err != nil {
					fmt.Println("Predict error:" + err.Error())
					disconnected = true
				}

				predictRunning = false
				requestCount++
			}()
		}
	}).ServeHTTP(w, req)
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&modelFNAME, "m", "./ggml-llama_7b-q4_1.bin", "path to quantized ggml model file to load")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	// flags.IntVar(&tokens, "n", 256, "number of tokens to predict")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(modelFNAME); os.IsNotExist(err) {
		fmt.Printf("Model file %s does not exist", modelFNAME)
		os.Exit(1)
	}

	listenURI := Address + ":" + Port
	uri := HttpProtocol + "://" + Address + ":" + Port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})
	http.HandleFunc("/ws", wsController)

	fmt.Printf("Server is running on %s\n\n", uri)

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
