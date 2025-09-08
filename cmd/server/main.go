package main

import (
	"log"
	"gpt-router/internal/server"
)

func main() {
	log.Println("🚀 GPT Router service starting...")
	server.Start()
}
