package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	prompt "example/focus-api/prompts"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type RequestBody struct {
	Goals []string `json:"goals" binding:"required"`
}

type ResponseBody struct {
	Phrases []string `json:"phrases"`
}

func generatePhrases(clientOpenAi *openai.Client, goals []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	ch := make(chan string, 100)
	var phrases []string

	for _, goal := range goals {
		wg.Add(1)
		go func(goal string) {
			defer wg.Done()

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
							Content: prompt.HumanMessage(goal),
						},
					},
				},
			)

			if err != nil {
				log.Printf("failed to generate phrases for goal %s: %v", goal, err)
				return
			}

			for _, choice := range resp.Choices {
				ch <- choice.Message.Content
			}
		}(goal)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		phrases = append(phrases, v)
	}

	if len(phrases) == 0 {
		return nil, errors.New("no phrases were generated")
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
		// This will triple the goals
		// goals := append(reqBody.Goals, append(reqBody.Goals, reqBody.Goals...)...)
		phrases, err := generatePhrases(clientOpenAi, reqBody.Goals) //  Use the goals variable to triple the requests per goal
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ResponseBody{Phrases: phrases})
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_KEY")
	clientOpenAi := openai.NewClient(apiKey)

	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())

	router.POST("/frases", handleGeneratePhrases(clientOpenAi))

	port := os.Getenv("PORT") // use 8080 in a non-production environment
	router.Run(":" + port)
}
