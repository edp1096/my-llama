package main // import "rest-server"

func main() {
	var err error
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"

	Llm, err = setupLLM(modelName)
	if err != nil {
		panic(err)
	}

	app := setupRouter()

	app.Listen("127.0.0.1:8864")
}
