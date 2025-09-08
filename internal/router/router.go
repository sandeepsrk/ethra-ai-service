package router

import (
	"context"
	"log"
	"os"
	"strings"

	openai "github.com/openai/openai-go"
)

var client *openai.Client

func init() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ OPENAI_API_KEY not set")
	}
	client = openai.NewClient(apiKey)
}

func RoutePrompt(ctx context.Context, prompt string) (string, error) {
	systemPrompt := `You are a router agent. Classify the user prompt into one of these categories:
- coding
- finance
- travel
- general
Respond ONLY with the category name.`

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4oMini),
		Messages: openai.F([]openai.ChatCompletionMessageParam{
			{
				Role: openai.F("system"),
				Content: openai.F([]openai.ChatCompletionContentPartParam{
					openai.ChatCompletionContentPartTextParam{
						Type: openai.F("text"),
						Text: systemPrompt,
					},
				}),
			},
			{
				Role: openai.F("user"),
				Content: openai.F([]openai.ChatCompletionContentPartParam{
					openai.ChatCompletionContentPartTextParam{
						Type: openai.F("text"),
						Text: prompt,
					},
				}),
			},
		}),
	})
	if err != nil {
		return "", err
	}

	category := strings.ToLower(strings.TrimSpace(resp.Choices[0].Message.Content[0].Text))
	return category, nil
}
