# ethra-ai-service

## install dependencies

go mod tidy

## set API key

export OPENAI_API_KEY="your_openai_api_key_here"

## run the service

go run ./cmd/server

## Then test:

curl -X POST http://localhost:5000/ask \
 -H "Content-Type: application/json" \
 -d '{"prompt": "Explain recursion in Go"}'
