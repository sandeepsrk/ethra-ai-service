package imageagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"ethra-go/internal/taxagent"
	"ethra-go/internal/types"

	openai "github.com/sashabaranov/go-openai"
	"github.com/joho/godotenv"
)

var client *openai.Client

func init() {
	_ = godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set in ImageAgent")
	}
	client = openai.NewClient(apiKey)
}

// ProcessImage takes base64 image and returns normalized JSON
func ProcessImage(base64Data string) (types.AgentResponse, error) {
	ctx := context.Background()

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "user",
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: `You are a bill analyzer. Extract all details from this image as JSON:
{
  "invoice_details": { "invoice_number": string, "date": string, "vendor": string, "address": string },
  "items": [ { "name": string, "quantity": number, "price": number, "total": number } ],
  "taxes": [ { "rate": number, "amount": number } ],
  "totals": { "subtotal": number, "discount": number, "tax_total": number, "grand_total": number, "round_off": number }
}`,
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: fmt.Sprintf("data:image/png;base64,%s", base64Data),
						},
					},
				},
			},
		},
		MaxTokens: 1500,
	})
	if err != nil {
		return types.AgentResponse{
			Agent:   "ImageAgent",
			Message: map[string]interface{}{"error": err.Error()},
		}, err
	}

	content := resp.Choices[0].Message.Content
	log.Printf("🧾 [ImageAgent] GPT raw output: %.300s", content)

	var extracted map[string]interface{}
	if err := json.Unmarshal([]byte(content), &extracted); err != nil {
		// If GPT returned text, keep raw
		extracted = map[string]interface{}{"raw": content}
	}

	// Normalize via TaxAgent
	taxResponse := taxagent.ProcessTaxGPT(extracted)

	final := types.AgentResponse{
		Agent: "Image+TaxAgent",
		Message: map[string]interface{}{
			"extracted":  extracted,
			"normalized": taxResponse.Message,
		},
	}
	return final, nil
}
