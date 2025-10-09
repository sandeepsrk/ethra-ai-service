package router

import (
	"net/http"
	
	"ethra-go/internal/handlers"
	"ethra-go/internal/middleware"

)

func RegisterRoutes() {
	mux := http.NewServeMux()
	mux.Handle("/prompt", middleware.APIKeyAuth(http.HandlerFunc(handlers.PromptHandler)))
}