package cmd

import (
	"fmt"
	"globo_test/cmd/api"
	"globo_test/cmd/incrementtest"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	rootCmd.AddCommand(api.ApiCommand())
	rootCmd.AddCommand(api.QueryApiCommand())
	rootCmd.AddCommand(api.CommandApiCommand())
	rootCmd.AddCommand(incrementtest.IncrementTestCommand())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
