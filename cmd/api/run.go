package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sergiodii/bbb/cmd/api/middleware"
	"github.com/spf13/cobra"
)

func ApiCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "api",
		Short: "Inicia a API de comandos e consultas",
	}

	c.Flags().StringP("port", "p", "8080", "Porta que a API irá escutar")
	c.Run = func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		fmt.Printf("\n[STARTING API] Iniciando API de comandos na porta %s...\n", port)
		r := gin.Default()

		// This middleware simulate the blocking of IP ranges
		// You can set the environment variable BLOCKED_IP_RANGES to a comma-separated list of IP prefixes to block
		// Example: export BLOCKED_IP_RANGES="192.168.1.,10.0.0."
		// Any request from an IP starting with these prefixes will be blocked with a 403 response
		r.Use(middleware.NewBlockingIPRangeMiddlewareV1())

		// This middleware simulate a rate limiting of 60 requests per minute per IP
		// In a real implementation, you would use a more robust solution with a datastore or in-memory structure
		// to track requests per IP and enforce limits.
		// Here, for simplicity, we just allow all requests.
		r.Use(middleware.RateLimitMiddlewareV1())

		queryApiRegister(r, "/query")
		commandApiRegister(r, "/command")
		r.Run(":" + port)
	}

	return &c
}

func QueryApiCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "query-api",
		Short: "Inicia a API de consultas",
	}

	c.Flags().StringP("port", "p", "8081", "Porta que a API irá escutar")
	c.Run = func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		fmt.Printf("\n[STARTING QUERY-API] Iniciando API de consultas na porta %s...\n", port)
		r := gin.Default()
		queryApiRegister(r, "")
		r.Run(":" + port)
	}

	return &c
}

func CommandApiCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "command-api",
		Short: "Inicia a API de comandos",
	}

	c.Flags().StringP("port", "p", "8082", "Porta que a API irá escutar")
	c.Run = func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		fmt.Printf("\n[STARTING COMMAND-API] Iniciando API de comandos na porta %s...\n", port)
		r := gin.Default()
		commandApiRegister(r, "")
		r.Run(":" + port)
	}

	return &c
}
