package agents

import (
	"encoding/base64"
	"log"
	"strings"

	"ethra-go/internal/imageagent"
	"ethra-go/internal/taxagent"
	"ethra-go/internal/types"
)

// RouteRequest decides which agent to use
func RouteRequest(agentType, content string, file *types.FileData) types.AgentResponse {
	switch strings.ToLower(agentType) {

	case "image":
		if file != nil && len(file.Buffer) > 0 {
			base64Data := base64.StdEncoding.EncodeToString(file.Buffer)
			log.Printf("🔹 [Agents] Sending image to ImageAgent (%d bytes)", len(file.Buffer))
			resp, err := imageagent.ProcessImage(base64Data)
			if err != nil {
				return types.AgentResponse{
					Agent:   "ImageAgent",
					Message: map[string]interface{}{"error": err.Error()},
				}
			}
			return resp
		}
		return types.AgentResponse{
			Agent:   "ImageAgent",
			Message: map[string]interface{}{"error": "No image data"},
		}

	case "tax":
		log.Println("🔹 [Agents] Sending text to TaxAgent")
		return taxagent.ProcessTaxGPT(content)

	case "auto":
		if file != nil && len(file.Buffer) > 0 {
			base64Data := base64.StdEncoding.EncodeToString(file.Buffer)
			resp, err := imageagent.ProcessImage(base64Data)
			if err != nil {
				return types.AgentResponse{
					Agent:   "ImageAgent",
					Message: map[string]interface{}{"error": err.Error()},
				}
			}
			return resp
		}
		return taxagent.ProcessTaxGPT(content)

	default:
		return types.AgentResponse{
			Agent:   "Router",
			Message: map[string]interface{}{"error": "Unknown agent type"},
		}
	}
}
