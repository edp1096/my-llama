package main

import (
	llama "github.com/edp1096/my-llama"
	"github.com/gofiber/fiber/v2"
)

func setupLLM(modelName string) (*llama.LLama, error) {
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

	app.Post("/v1/chat/completions", chatCompletions)

	return app
}
