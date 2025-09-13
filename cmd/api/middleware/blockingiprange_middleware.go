package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewBlockingIPRangeMiddlewareV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := getClientIP(c)

		blockedRanges := getBlockedIPRanges()
		for _, r := range blockedRanges {
			if strings.HasPrefix(clientIP, strings.TrimSpace(r)) {
				// Bloqueia o acesso
				c.AbortWithStatusJSON(403, gin.H{"error": "Access from your IP range is blocked"})
				return
			}
		}

		c.Next()
	}
}

func getBlockedIPRanges() []string {
	r := os.Getenv("BLOCKED_IP_RANGES")
	if r == "" {
		return []string{}
	}
	return strings.Split(r, ",")
}
