package main // import "run-myllama"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	_ "embed"

	"github.com/shirou/gopsutil/v3/cpu"
	ws "golang.org/x/net/websocket"

	llama "github.com/edp1096/my-llama"
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
	threads  = 4
	useMlock = true
)

var (
	cpuPhysicalNUM = 0
	cpuLogicalNUM  = 0

	isBrowserOpen = false

	modelPath = "./"

	modelFname  string   = ""
	modelFnames []string = []string{}
)

func setQueryParams(l *llama.LLama, req *http.Request) {
	// Model file
	model_file := req.URL.Query().Get("model_file")
	for _, fname := range modelFnames {
		if strings.TrimSpace(fname) == strings.TrimSpace(model_file) {
			modelFname = model_file
			break
		}
	}

	nThreadsSTR := req.URL.Query().Get("threads")
	if nThreadsSTR != "" {
		nThreads, err := strconv.Atoi(nThreadsSTR)
		if err != nil {
			l.Threads = nThreads
		}
	}

	useDumpSessionSTR := req.URL.Query().Get("use_dump_session")
	if useDumpSessionSTR != "" {
		l.UseDumpSession = false
		if useDumpSessionSTR == "true" {
			l.UseDumpSession = true
		}
	}

	nCtxSTR := req.URL.Query().Get("n_ctx")
	if nCtxSTR != "" {
		nCTX, err := strconv.Atoi(nCtxSTR)
		if err == nil {
			l.SetNCtx(nCTX)
		}
	}

	nBatchSTR := req.URL.Query().Get("n_batch")
	if nBatchSTR != "" {
		nBATCH, err := strconv.Atoi(nBatchSTR)
		if err != nil {
			l.SetNBatch(nBATCH)
		}
	}

}

func evalAndResponse(l *llama.LLama, conn *ws.Conn, handler ws.Codec) error {
	remainCOUNT := l.GetRemainCount()

	responseBufferBytes := []byte{}
	responseBuffer := ""

END:
	for remainCOUNT != 0 {
		ok := l.PredictTokens()
		if !ok {
			return fmt.Errorf("failed to predict the tokens")
		}

		// display text
		embdSIZE := l.GetEmbedSize()
		for i := 0; i < embdSIZE; i++ {
			select {
			case <-l.PredictStop:
				remainCOUNT = 0
				break END
			default:
				embedSTR := l.GetEmbedString(i)

				responseBufferBytes = append(responseBufferBytes, []byte(embedSTR)...)
				if !utf8.ValidString(embedSTR) {
					continue // Because connection is closed, don't send invalid UTF-8
				}

				if len(responseBufferBytes) > 0 {
					responseBuffer = string(responseBufferBytes)
					if !utf8.ValidString(responseBuffer) {
						continue
					}
				} else {
					responseBuffer += embedSTR
				}

				// fmt.Print(responseBuffer)

				err := handler.Send(conn, "$$__RESPONSE_PREDICT__$$\n$$__SEPARATOR__$$\n"+responseBuffer)
				if err != nil {
					fmt.Println("Send error:", err)
					remainCOUNT = 0
					break END
				}

				responseBufferBytes = []byte{}
				responseBuffer = ""
			}
		}

		ok = l.CheckPromptOrContinue()
		if !ok {
			break
		}
		l.DropBackUserInput()

		remainCOUNT = l.GetRemainCount()
	}

	err := handler.Send(conn, "$$__RESPONSE_PREDICT__$$\n$$__SEPARATOR__$$\n"+"\n$$__RESPONSE_DONE__$$\n")
	if err != nil {
		fmt.Println("Send error:", err)
		return err
	}

	l.PrintTimings()

	return nil
}

func wsController(w http.ResponseWriter, req *http.Request) {
	ws.Handler(func(conn *ws.Conn) {
		l, err := llama.New()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer l.FreeALL()
		defer conn.Close()

		l.Threads = threads

		req := conn.Request()
		setQueryParams(l, req)

		requestCount := 0 // to ignore prompt

		l.InitParams()

		l.SetThreadsCount(l.Threads)
		l.SetUseMlock(useMlock)

		fmt.Println("Threads:", l.Threads)

		err = l.LoadModel(modelFname)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Model initialized..")

		dumpFname := `dumpsession_` + modelFname + `.session`
		dumpInitialLoaded := false

		// Load dump_session
		if l.UseDumpSession {
			if _, err := os.Stat(dumpFname); err == nil {
				fmt.Println("Load", dumpFname)
				l.LoadSession(dumpFname)
				dumpInitialLoaded = true
			}
		}

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
				switch true {
				case strings.HasPrefix(message, "$$__COMMAND__$$"):
					command := strings.Split(message, "\n$$__SEPARATOR__$$\n")[1]

					switch command {
					case "$$__KILL_SERVER__$$":
						os.Exit(0)
					case "$$__STOP__$$":
						if predictRunning {
							l.PredictStop <- true
						}
					case "$$__DEVICE_TYPE__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__DEVICE_TYPE__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%s", tag, deviceType)
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
					case "$$__DUMPSESSION_EXIST__$$":
						// Check dumpsession file exists
						if _, err := os.Stat(dumpFname); err == nil {
							err = handler.Send(conn, "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__DUMPSESSION_EXIST__$$\n$$__SEPARATOR__$$\ntrue")
							if err != nil {
								fmt.Println("Send error:", err)
								disconnected = true
							}
						} else {
							err = handler.Send(conn, "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__DUMPSESSION_EXIST__$$\n$$__SEPARATOR__$$\nfalse")
							if err != nil {
								fmt.Println("Send error:", err)
								disconnected = true
							}
						}
					case "$$__MODEL_FILE__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__MODEL_FILES__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%s", tag, strings.Join(modelFnames, "\n$$__SEPARATOR__$$\n"))
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
					case "$$__MODEL_FILES__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__MODEL_FILES__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%s", tag, strings.Join(modelFnames, "\n$$__SEPARATOR__$$\n"))
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
					case "$$__MAX_CPU_PHYSICAL__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_PHYSICAL__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%d", tag, cpuPhysicalNUM)
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
					case "$$__MAX_CPU_LOGICAL__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_LOGICAL__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%d", tag, cpuLogicalNUM)
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
					case "$$__THREADS__$$":
						tag := "$$__RESPONSE_INFO__$$\n$$__SEPARATOR__$$\n$$__THREADS__$$\n$$__SEPARATOR__$$\n"
						response := fmt.Sprintf("%s%d", tag, l.Threads)
						err = handler.Send(conn, response)
						if err != nil {
							fmt.Println("Send error:", err)
							disconnected = true
						}
					}

					continue
				case strings.HasPrefix(message, "$$__PARAMETER__$$"):
					params := strings.Split(message, "\n$$__SEPARATOR__$$\n")
					paramNAME := params[1]
					paramVALUE := strings.TrimSpace(params[2])

					switch paramNAME {
					case "$$__THREADS__$$":
						// Maybe not needed
						threadsReceived, _ := strconv.Atoi(paramVALUE)
						fmt.Println("threads:", threadsReceived)
						if threadsReceived > 3 {
							threads = threadsReceived
						}
					case "$$__N_CTX__$$":
						n_ctx, _ := strconv.Atoi(paramVALUE)
						fmt.Println("n_ctx:", n_ctx)
						l.SetNCtx(n_ctx)
					case "$$__N_BATCH__$$":
						n_batch, _ := strconv.Atoi(paramVALUE)
						fmt.Println("n_batch:", n_batch)
						l.SetNBatch(n_batch)
					case "$$__SAMPLING_METHOD__$$":
						samplingMethod := paramVALUE
						fmt.Println("Sampling method:", samplingMethod)
						l.SetSamplingMethod(samplingMethod)
					case "$$__TOP_K__$$":
						top_k, _ := strconv.Atoi(paramVALUE)
						fmt.Println("top_k:", top_k)
						l.SetTopK(top_k)
					case "$$__TOP_P__$$":
						top_p, _ := strconv.ParseFloat(paramVALUE, 64)
						fmt.Println("top_p:", top_p)
						l.SetTopP(top_p)
					case "$$__TEMPERATURE__$$":
						temperature, _ := strconv.ParseFloat(paramVALUE, 64)
						fmt.Println("temperature:", temperature)
						l.SetTemperature(temperature)
					case "$$__REPEAT_PENALTY__$$":
						penalty, _ := strconv.ParseFloat(paramVALUE, 64)
						fmt.Println("repeat penalty:", penalty)
						l.SetRepeatPenalty(penalty)
					default:
						fmt.Println("Unknown parameter:", paramNAME)
					}

					continue
				}

				// fmt.Printf("%s\n", message) // Print received message from client
			}

			if predictRunning {
				continue
			}

			datas := strings.Split(message, "\n$$__SEPARATOR__$$\n")

			// Predict and Write responses
			go func() {
				// runtime.LockOSThread()

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
					if dumpInitialLoaded {
						// Because the prompt is not needed after the reload dump_session
						l.SetPrompt(antiprompt)
					} else {
						l.SetPrompt(reflectionPrompt)
					}
					l.SetAntiPrompt(antiprompt)

					err = l.AllocateVariables()
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

				// err = l.Predict(conn, handler)
				err = evalAndResponse(l, conn, handler)
				if err != nil {
					fmt.Println("evalAndResponse error:" + err.Error())
					disconnected = true
				}

				// Save dump_session
				if l.UseDumpSession {
					fmt.Println("Save", dumpFname)
					l.SaveSession(dumpFname)
				}

				predictRunning = false
				requestCount++
			}()
		}

		for predictRunning {
			time.Sleep(100 * time.Millisecond)
		}

	}).ServeHTTP(w, req)
}

func main() {
	cpuPhysicalNUM, _ = cpu.Counts(false)
	cpuLogicalNUM, _ = cpu.Counts(true)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.BoolVar(&isBrowserOpen, "b", false, "open browser automatically")

	threads = cpuPhysicalNUM

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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, index_html)
	})
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		fmt.Fprint(w, style_css)
	})
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		fmt.Fprint(w, script_js)
	})
	http.HandleFunc("/ws", wsController)

	fmt.Printf("Server is running on %s\n\n", uri)

	if isBrowserOpen {
		openBrowser(uri)
	}

	log.Fatal(http.ListenAndServe(listenURI, nil))
}
