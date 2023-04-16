package main // import "my-llama"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "embed"

	"github.com/shirou/gopsutil/v3/cpu"
	ws "golang.org/x/net/websocket"

	llama "my-llama/cgollama"
)

//go:embed html/index.html
var index_html string

//go:embed html/style.css
var style_css string

//go:embed html/script.js
var script_js string

var (
	Address      string = "localhost"
	Port         string = "1323"
	HttpProtocol string = "http"
)

var (
	threads = 4
	// tokens  = 256
)

var (
	isBrowserOpen = false

	modelPath = "./"

	modelFname  string   = ""
	modelFnames []string = []string{}
)

func wsController(w http.ResponseWriter, req *http.Request) {
	ws.Handler(func(conn *ws.Conn) {
		defer conn.Close()

		req := conn.Request()
		model_file := req.URL.Query().Get("model_file")
		if model_file != "" {
			// modelFNAME = model_file
			fmt.Println("model_file:", model_file)
		}

		topk := req.URL.Query().Get("topk")
		if topk != "" {
			// tokens, _ = strconv.Atoi(topk)
			fmt.Println("topk:", topk)
		} else {
			// tokens = 256
			fmt.Println("topk is not set")
		}

		var err error

		requestCount := 0 // to ignore prompt

		l, err := llama.New()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = l.LoadModel(modelFname)
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
	cpuCoreNUM, _ := cpu.Counts(false)
	cpuLogicalNUM, _ := cpu.Counts(true)
	fmt.Println("CPU cores/logical:", cpuCoreNUM, "/", cpuLogicalNUM)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.BoolVar(&isBrowserOpen, "b", false, "open browser automatically")
	flags.StringVar(&modelFname, "m", "", "path to quantized ggml model file to load")
	flags.IntVar(&threads, "t", cpuCoreNUM, "number of threads to use during computation")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	if modelFname == "" {
		modelFnames, err = findModelFiles(modelPath)
		if err != nil {
			fmt.Printf("Finding model files failed: %s", err)
			os.Exit(1)
		}
		// fmt.Println(modelFnames)

		if len(modelFnames) == 0 {
			fmt.Println("No model files found. Download the model file(s) before launch.")
			fmt.Println("Press enter to open the model search page.")
			fmt.Scanln()

			openBrowser("https://huggingface.co/search/full-text?q=ggml+7b&type=model")
			os.Exit(1)
		}

		modelFname = modelFnames[0]
	}

	if _, err := os.Stat(modelFname); os.IsNotExist(err) {
		fmt.Printf("Model file %s does not exist", modelFname)
		os.Exit(1)
	}

	listenURI := Address + ":" + Port
	uri := HttpProtocol + "://" + Address + ":" + Port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, index_html)
	})
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, style_css)
	})
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, script_js)
	})
	http.HandleFunc("/ws", wsController)

	fmt.Printf("Server is running on %s\n\n", uri)

	if isBrowserOpen {
		openBrowser(uri)
	}

	log.Fatal(http.ListenAndServe(listenURI, nil))
}
