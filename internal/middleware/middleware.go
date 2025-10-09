package middleware

import (
	"net/http"
	"os"
    "log"
	"encoding/json"

	"ethra-go/internal/types"

	"github.com/joho/godotenv"

)

// APIKeyAuthMiddleware checks for an "X-API-Key" header
func APIKeyAuth(next http.Handler) http.Handler {
	_ = godotenv.Load()
	apiKey := os.Getenv("NODE_API_KEY")
	if apiKey == "" {
		log.Fatal("NODE_API_KEY not set in ImageAgent")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqKey := r.Header.Get("X-API-Key")
		if reqKey == "" || reqKey != apiKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(types.ErrorResponse{
				Error: "unauthorized, invalid API key",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}