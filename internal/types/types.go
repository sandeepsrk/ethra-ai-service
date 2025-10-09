package types

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type AgentResponse struct {
	Agent   string      `json:"agent"`
	Message interface{} `json:"message"`
}

type FileData struct {
	Buffer []byte
	Name   string
	Type   string
}

type MemoryEntry struct {
	Role    string
	Content interface{}
	Agent   string
}


type ErrorResponse struct {
	Error string `json:"error"`
}



