# ethra-ai-service

This is a Go-based multi-agent router service powered by OpenAI GPT models.  
It listens for prompts, classifies them into categories (coding, finance, travel, general),  
and dispatches them to specialized agents.

---

## Install dependencies

```bash
go mod tidy
```

## Set API key

```bash
export OPENAI_API_KEY="your_openai_api_key_here"
```

## Run the service

```bash
go run ./cmd/server
```

## Then test:

```bash
curl -X POST http://localhost:5000/ask  -H "Content-Type: application/json"  -d '{"prompt": "Explain recursion in Go"}'
```

---

## Project Structure

```
ethra-ai-service/
 ├── go.mod
 ├── cmd/server/main.go        # Entry point
 └── internal/
     ├── server/server.go      # HTTP server setup
     ├── router/router.go      # Router agent (OpenAI based)
     └── agents/agents.go      # Specialized agents
```
