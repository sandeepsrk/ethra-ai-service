package router

import (
	"net/http"
	
	"ethra-go/internal/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/prompt", handlers.PromptHandler)
}