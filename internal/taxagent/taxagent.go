package taxagent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"reflect"

	"ethra-go/internal/types"
	
	openai "github.com/sashabaranov/go-openai"
	"github.com/joho/godotenv"
)

var client *openai.Client

func init() {
	_ = godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set in TaxAgent")
	}
	client = openai.NewClient(apiKey)
}

// ProcessTaxGPT handles text or JSON input and returns normalized invoice JSON.
// Any non-JSON GPT output will be wrapped in {"raw": "..."} to ensure valid JSON.
func ProcessTaxGPT(input interface{}) types.AgentResponse {
	start := time.Now()
	log.Println("🧮 [TaxAgent] Starting normalization...")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var textInput string
	switch v := input.(type) {
	case string:
		textInput = v
		log.Printf("📜 [TaxAgent] Received text input: %.200s...", v)
	case map[string]interface{}:
		b, _ := json.MarshalIndent(v, "", "  ")
		textInput = string(b)
		log.Printf("🧾 [TaxAgent] Received JSON input: %.300s...", textInput)
	default:
		textInput = fmt.Sprintf("%v", v)
		log.Printf("⚠️ [TaxAgent] Unknown input type (%s) — coerced to string", reflect.TypeOf(input).String())
	}

	systemMsg := `You are a tax assistant. Normalize extracted invoice data into this JSON schema:
{
  "invoice_details": { "invoice_number": string, "date": string, "vendor": string, "address": string },
  "items": [ { "name": string, "quantity": number, "price": number, "total": number, tax: { "rate": number, "amount": number } } ],
  "totals": { "subtotal": number, "discount": number, "tax_total": number, "grand_total": number, "round_off": number }
}
Always return valid JSON.`

	log.Println("🔹 [TaxAgent] Sending request to GPT...")

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: systemMsg,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Normalize this bill:\n%s", textInput),
			},
		},
		MaxTokens: 1500,
	})
	if err != nil {
		log.Printf("❌ [TaxAgent] GPT error: %v", err)
		return types.AgentResponse{
			Agent:   "TaxAgent",
			Message: map[string]interface{}{"error": err.Error()},
		}
	}

	if len(resp.Choices) == 0 {
		log.Printf("⚠️ [TaxAgent] Empty GPT response")
		return types.AgentResponse{
			Agent:   "TaxAgent",
			Message: map[string]interface{}{"raw": ""},
		}
	}

	output := strings.TrimSpace(resp.Choices[0].Message.Content)
	log.Printf("🧩 [TaxAgent] GPT Output: %.300s...", output)

	var normalized map[string]interface{}
	if err := json.Unmarshal([]byte(output), &normalized); err != nil {
		log.Printf("⚠️ [TaxAgent] Invalid JSON — wrapping raw text")
		normalized = map[string]interface{}{"raw": output}
	}

	log.Printf("✅ [TaxAgent] Normalization complete in %v", time.Since(start))
	return types.AgentResponse{Agent: "TaxAgent", Message: normalized}
}
