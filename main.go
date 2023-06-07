package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	prompt "example/focus-api/prompts"
	q "example/focus-api/queue"

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

func generatePhrases(apiKeys []string, goals []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 270*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	ch := make(chan string, 100)
	var phrases []string

	for index, goal := range goals {
		wg.Add(1)
		clientOpenAi := openai.NewClient(apiKeys[index])
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
		phrasesArray := strings.Split(v, ";")
		for i := range phrasesArray {
			phrasesArray[i] = strings.ReplaceAll(phrasesArray[i], `"`, "")
			phrasesArray[i] = strings.TrimSpace(phrasesArray[i])
		}
		phrases = append(phrases, phrasesArray...)
	}

	if len(phrases) == 0 {
		return nil, errors.New("no phrases were generated")
	}

	return phrases, nil
}

func handleGeneratePhrases(q *q.Queue) gin.HandlerFunc {
	return func(c *gin.Context) {

		var reqBody RequestBody
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(reqBody.Goals) > 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit of 6 goals per request exceeded"})
			return
		}

		goals := append(reqBody.Goals, reqBody.Goals...)

		apiKeys := make([]string, 30)

		for index := range goals {
			apiKey, _ := (*q).Dequeue()
			(*q).Enqueue(apiKey)
			apiKeys[index] = apiKey
		}

		phrases, err := generatePhrases(apiKeys, goals)
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

	keyQueue := q.Queue{
		os.Getenv("OPENAI_KEY1"),
		os.Getenv("OPENAI_KEY2"),
		os.Getenv("OPENAI_KEY3"),
		os.Getenv("OPENAI_KEY4"),
		os.Getenv("OPENAI_KEY5"),
		os.Getenv("OPENAI_KEY6"),
		os.Getenv("OPENAI_KEY7"),
		os.Getenv("OPENAI_KEY8"),
		os.Getenv("OPENAI_KEY9"),
		os.Getenv("OPENAI_KEY10"),
		os.Getenv("OPENAI_KEY11"),
		os.Getenv("OPENAI_KEY12"),
		os.Getenv("OPENAI_KEY13"),
		os.Getenv("OPENAI_KEY14"),
		os.Getenv("OPENAI_KEY15"),
		os.Getenv("OPENAI_KEY16"),
		os.Getenv("OPENAI_KEY17"),
		os.Getenv("OPENAI_KEY18"),
		os.Getenv("OPENAI_KEY19"),
		os.Getenv("OPENAI_KEY20"),
	}

	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())

	router.POST("/frases", handleGeneratePhrases(&keyQueue))

	port := os.Getenv("PORT") // use 8080 in a non-production environment
	router.Run(":" + port)
}
