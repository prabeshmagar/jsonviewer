package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [file] [query] [output]",
	Short: "Export specific JSON data to a new file",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		query := args[1]
		output := args[2]

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		result := gjson.GetBytes(data, query)
		if !result.Exists() {
			fmt.Println("Query did not match any data.")
			return
		}

		exportData, err := json.MarshalIndent(result.Value(), "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		if err := os.WriteFile(output, exportData, 0644); err != nil {
			fmt.Println("Error writing file:", err)
			return
		}

		fmt.Printf("Exported data to %s\n", output)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

