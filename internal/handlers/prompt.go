package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"

	"ethra-go/internal/agents"
	"ethra-go/internal/types"
)

// Handle POST /prompt
func PromptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(20 << 20) // 20 MB max
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	prompt := r.FormValue("prompt")
	sessionId := r.FormValue("sessionId")

	var fileData *types.FileData
	file, header, err := r.FormFile("file")
	if err == nil && file != nil {
		defer file.Close()
		buf, _ := ioutil.ReadAll(file)
		fileData = &types.FileData{
			Buffer: buf,
			Name:   header.Filename,
			Type:   header.Header.Get("Content-Type"),
		}
		log.Printf("📎 [PromptHandler] Received file %s (%d bytes)", header.Filename, len(buf))
	}

	// Decide agent type: image if file present, auto otherwise
	agentType := "auto"
	if fileData != nil {
		agentType = "image"
	}

	log.Printf("📥 [PromptHandler] sessionId=%s, agentType=%s, prompt=%s", sessionId, agentType, prompt)

	resp := agents.RouteRequest(agentType, prompt, fileData)

	log.Printf("📝 [PromptHandler] Agent response: %+v", resp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
