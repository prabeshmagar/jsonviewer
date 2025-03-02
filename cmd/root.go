package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jsonviewer",
	Short: "A CLI JSON viewer with tree and search support",
	Long:  "jsonviewer is a tool for viewing and navigating JSON files in a terminal.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

