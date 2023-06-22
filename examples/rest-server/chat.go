package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model,omitempty"`
	Messages []ChatMessage `json:"messages"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Choices []ChatChoice `json:"choices"`
	Usage   ChatUsage    `json:"usage"`
}

func setupRouter() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/v1/chat/hello", helloChat)
	app.Post("/v1/chat/completions", chatCompletions)

	return app
}

func helloChat(c *fiber.Ctx) error {
	creq := new(ChatRequest)
	c.BodyParser(creq)

	fmt.Println(creq)

	response := new(ChatResponse)
	response.ID = "test"
	response.Object = "chat.completion"
	response.Created = time.Now().Unix()

	choice := new(ChatChoice)
	choice.Index = 0
	choice.Message = ChatMessage{Role: "bot", Content: "Hello, World 👋!"}
	choice.FinishReason = "stop"

	response.Choices = append(response.Choices, *choice)
	response.Usage.PromptTokens = 2
	response.Usage.CompletionTokens = 7
	response.Usage.TotalTokens = response.Usage.PromptTokens + response.Usage.CompletionTokens

	return c.JSON(response)
}

func chatCompletions(c *fiber.Ctx) error {
	var err error

	creq := new(ChatRequest)
	c.BodyParser(creq)

	messageContent := ""

	response := new(ChatResponse)
	response.ID = "test"
	response.Object = "chat.completion"
	response.Created = time.Now().Unix()

	numPromptTokens := 0
	numCompletionTokens := 0
	numTotalTokens := 0

	numPast := 0

	messages := creq.Messages
	prompt := ""
	for _, message := range messages {
		prompt += message.Role + ": " + message.Content
	}

	numPast, err = evalPrompt(prompt, numPast)
	if err != nil {
		return err
	}
	numPromptTokens = numPast

	lineBreakCount := 0
	for i := 0; i < 16; i++ {
		// i := 0
		// for {
		var nextTokenStr string
		nextTokenStr, numPast = getTokenString(numPast)
		messageContent += nextTokenStr

		if nextTokenStr == "\n" {
			lineBreakCount++
			if lineBreakCount >= 2 {
				break
			}
		} else {
			lineBreakCount = 0
		}

		// i++
	}

	numCompletionTokens = numPast - numPromptTokens
	numTotalTokens = numPast

	choice := &ChatChoice{
		Index:        0,
		Message:      ChatMessage{Role: "bot", Content: messageContent},
		FinishReason: "stop",
	}
	response.Choices = append(response.Choices, *choice)

	response.Usage = ChatUsage{
		PromptTokens:     numPromptTokens,
		CompletionTokens: numCompletionTokens,
		TotalTokens:      numTotalTokens,
	}

	return c.JSON(response)
}
