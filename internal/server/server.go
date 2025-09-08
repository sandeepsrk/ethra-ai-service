package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gpt-router/internal/agents"
	"gpt-router/internal/router"
)

func Start() {
	r := gin.Default()

	r.POST("/ask", func(c *gin.Context) {
		var body struct {
			Prompt string `json:"prompt"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		category, err := router.RoutePrompt(c, body.Prompt)
		if err != nil {
			log.Println("Router error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Routing failed"})
			return
		}

		var result string
		switch category {
		case "coding":
			result = agents.CodingAgent(body.Prompt)
		case "finance":
			result = agents.FinanceAgent(body.Prompt)
		case "travel":
			result = agents.TravelAgent(body.Prompt)
		default:
			result = agents.GeneralAgent(body.Prompt)
		}

		c.JSON(http.StatusOK, gin.H{
			"category": category,
			"result":   result,
		})
	})

	r.Run(":5000")
}
