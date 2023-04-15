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
	tokens  = 256
)

var modelFNAME string

func initModel(modelFNAME string) *llama.LLama {
	l, err := llama.New(modelFNAME)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return l
}

func wsController(w http.ResponseWriter, req *http.Request) {
	ws.Handler(func(conn *ws.Conn) {
		defer conn.Close()

		l := initModel(modelFNAME)

		err := l.LoadModel(modelFNAME)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = l.InitParams()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		l.SetupParams()

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

			// Predict and Write
			go func() {
				predictRunning = true

				if len(datas) < 3 {
					return
				}

				prompt := datas[1]
				input := prompt + datas[0]
				// antiprompt := datas[2]

				// fmt.Println("Input:", "'"+input+"'")
				// fmt.Println("Prompt:", "'"+prompt+"'")
				// fmt.Println("Antiprompt:", "'"+antiprompt+"'")

				err := l.Predict(conn, handler, input)
				if err != nil {
					fmt.Println("Predict error:" + err.Error())
					disconnected = true
				}

				// params, predVARs, remainCOUNT = l.GetInitialParams(input, prompt, antiprompt, po)

				// err = l.Predict(conn, handler, input, prompt, antiprompt, params, predVARs, remainCOUNT)
				// if err != nil {
				// 	fmt.Println("Predict error:" + err.Error())
				// 	disconnected = true
				// }

				// predictRunning = false
				// requestCount++
				// fmt.Println("Request count:", requestCount)
			}()
		}
	}).ServeHTTP(w, req)
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&modelFNAME, "m", "./ggml-llama_7b-q4_1.bin", "path to quantized ggml model file to load")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	flags.IntVar(&tokens, "n", 256, "number of tokens to predict")

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
