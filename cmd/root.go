package cmd

import (
	"fmt"
	"os"

	"github.com/sergiodii/bbb/cmd/api"
	"github.com/sergiodii/bbb/cmd/loadtest"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	rootCmd.AddCommand(api.ApiCommand())
	rootCmd.AddCommand(api.QueryApiCommand())
	rootCmd.AddCommand(api.CommandApiCommand())
	rootCmd.AddCommand(loadtest.LoadTestCommand())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
