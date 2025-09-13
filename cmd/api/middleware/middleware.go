package middleware

import "github.com/gin-gonic/gin"

// isAllowed verifies if the IP can make a new request
// For simplicity, this example allows all requests.
func isAllowed(ip string) bool {

	// In a real implementation, we would check a datastore or in-memory structure
	// to count requests from this IP in the last minute and compare with the limit.

	// For this example, we will allow all requests.
	return true
}

// getClientIP extracts the real IP of the client (considering proxies)
func getClientIP(c *gin.Context) string {
	// Verifica headers de proxy comuns
	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		return ip // Cloudflare
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip // Nginx
	}
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return ip // Load balancers
	}

	// Fallback para IP direto
	return c.ClientIP()
}
