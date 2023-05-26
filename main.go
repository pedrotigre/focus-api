package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	prompt "example/focus-api/prompts"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type RequestBody struct {
	Topic string `json:"topic" binding:"required"`
}

type ResponseBody struct {
	Phrases []string `json:"phrases"`
}

func generatePhrases(clientOpenAi *openai.Client, topic string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := clientOpenAi.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,

			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt.SystemMessage(),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt.HumanMessage(topic),
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no phrases generated")
	}

	var phrases []string
	for _, choice := range resp.Choices {
		phrases = append(phrases, choice.Message.Content)
	}

	return phrases, nil
}

func handleGeneratePhrases(clientOpenAi *openai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody RequestBody
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		phrases, err := generatePhrases(clientOpenAi, reqBody.Topic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ResponseBody{Phrases: phrases})
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_KEY")
	clientOpenAi := openai.NewClient(apiKey)

	router := gin.Default()
	router.POST("/frases", handleGeneratePhrases(clientOpenAi))
	router.Run(":8080")
}
