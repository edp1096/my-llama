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
	cpuPhysicalNUM = 0
	cpuLogicalNUM  = 0

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

		// Todo: settings from query string
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

			// Command
			if len(message) > 0 {
				if strings.HasPrefix(message, "$$__COMMAND__$$") {
					command := strings.Split(message, "\n$$__SEPARATOR__$$\n")[1]

					switch command {
					case "$$__STOP__$$":
						if predictRunning {
							l.PredictStop <- true
						}
						continue
					case "$$__MAX_CPU_PHYSICAL__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_PHYSICAL__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%d", tag, cpuPhysicalNUM)
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
						continue
					case "$$__MAX_CPU_LOGICAL__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_PHYSICAL__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%d", tag, cpuLogicalNUM)
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
						continue
					}
				}

				// fmt.Printf("%s\n", message) // Print received message from client
			}

			if predictRunning {
				continue
			}

			datas := strings.Split(message, "\n$$__SEPARATOR__$$\n")

			// Predict and Write responses
			go func() {
				predictRunning = true

				if len(datas) < 3 {
					return
				}

				if datas[0] != "$$__PROMPT__$$" {
					return
				}

				// 0: prompt, 1: input, 2: reflection prompt, 3: antiprompt
				reflectionPrompt = datas[2]
				antiprompt = datas[3]
				input = datas[1]

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
	cpuPhysicalNUM, _ = cpu.Counts(false)
	cpuLogicalNUM, _ = cpu.Counts(true)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.BoolVar(&isBrowserOpen, "b", false, "open browser automatically")
	flags.StringVar(&modelFname, "m", "", "path to quantized ggml model file to load")
	flags.IntVar(&threads, "t", cpuPhysicalNUM, "number of threads to use during computation")

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
			fmt.Println("No model files found.")
			fmt.Println("Press enter to download vicuna model file and to open the model search page.")
			fmt.Println("Press Ctrl+C, if you want to exit.")
			fmt.Scanln()

			openBrowser(weightsSearchURL)
			downloadVicuna()

			modelFnames, _ = findModelFiles(modelPath)
		}

		modelFname = modelFnames[0]
	}

	if _, err := os.Stat(modelFname); os.IsNotExist(err) {
		// Because, model will be downloaded if not exists, may be not reachable here
		fmt.Printf("Model file %s does not exist", modelFname)
		os.Exit(1)
	}

	fmt.Println("CPU cores physical/logical:", cpuPhysicalNUM, "/", cpuLogicalNUM)

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
