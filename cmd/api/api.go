package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
