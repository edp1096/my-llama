package main // import "rest-server"

import (
	llama "github.com/edp1096/my-llama"
	"github.com/gofiber/fiber/v2"
)

func initLLM(modelName string) (*llama.LLama, error) {
	var err error
	numPredict := 16

	l, err := llama.New()
	if err != nil {
		return nil, err
	}

	l.LlamaApiInitBackend()
	l.InitGptParams()

	l.SetNumThreads(4)
	l.SetUseMlock(true)
	l.SetNumPredict(numPredict)
	l.SetNumGpuLayers(32)
	l.SetSeed(42)

	l.InitContextParamsFromGptParams()

	err = l.LoadModel(modelName)
	if err != nil {
		return nil, err
	}

	l.AllocateTokens()

	return l, err
}

func setupRouter() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/v1/chat/completions", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	return app
}

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	_, err := initLLM(modelName)
	if err != nil {
		panic(err)
	}

	app := setupRouter()

	app.Listen("127.0.0.1:8864")
}
