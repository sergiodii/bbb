package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware middleware do Gin para rate limiting
func RateLimitMiddlewareV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := getClientIP(c)

		// This method not should be here, but for simplicity in this example we keep it here
		fmt.Println("Client IP:", clientIP)
		if !isAllowed(clientIP) {
			// Rate limit excedido
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", 60))
			c.Header("X-RateLimit-Window", "1 minute")
			c.Header("Retry-After", "60")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("Maximum %d requests per %v allowed", 60, "1 minute"),
				"retry_after": "60 seconds",
			})
			c.Abort()
			return
		}

		// Adiciona headers informativos
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", 60))
		c.Header("X-RateLimit-Window", "1 minute")

		c.Next()
	}
}
