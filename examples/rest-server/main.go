package main // import "rest-server"

func main() {
	var err error

	// llmConfig := new(ModelExecutionConfig)
	// llmConfig.ModelName = "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	llmConfig := &ModelExecutionConfig{
		ModelName:    "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin",
		NumThreads:   4,
		UseMlock:     true,
		NumGpuLayers: 32,
		Seed:         42,
	}

	// samplingParams := new(SamplingParameters)
	// samplingParams.NumPredict = 16
	samplingParams := &SamplingParameters{
		NumPredict: 16,
	}

	err = setupLLM(*llmConfig, *samplingParams)
	if err != nil {
		panic(err)
	}

	app := setupRouter()
	app.Listen("127.0.0.1:8864")
}
