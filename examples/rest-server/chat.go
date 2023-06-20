package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func chatCompletions(c *fiber.Ctx) error {
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
