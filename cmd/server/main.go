package main

import (
	"log"
	"os"
	"net/http"

	"ethra-go/internal/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY not set in environment")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.RegisterRoutes()
	log.Println("🚀 Server running on http://localhost:"+port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
