package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var outputFile string

var searchCmd = &cobra.Command{
	Use:   "search [file] [key=value]",
	Short: "Search for objects containing a specific key-value pair in a JSON file.",
	Long: `Searches for JSON objects containing a specific key-value pair in the given JSON file. 
For example:
  search data.json name=Alice --output results.json`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		search := args[1]

		// Split the search argument into key and value
		parts := strings.SplitN(search, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Error: Search query must be in the format key=value.")
			return
		}
		key := parts[0]
		value := parts[1]

		// Read the JSON file
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			return
		}

		// Parse the JSON into a map for easier traversal
		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			fmt.Printf("Error parsing JSON file: %v\n", err)
			return
		}

		// Recursively search for matching objects
		matches := searchObjects(jsonData, key, value)

		if len(matches) == 0 {
			fmt.Printf("No objects found containing key=%s and value=%s.\n", key, value)
			return
		}

		// Either output to file or print results
		if outputFile != "" {
			fileData, err := json.MarshalIndent(matches, "", "  ")
			if err != nil {
				fmt.Printf("Error exporting results to file: %v\n", err)
				return
			}
			if err := os.WriteFile(outputFile, fileData, 0644); err != nil {
				fmt.Printf("Error writing to file %s: %v\n", outputFile, err)
				return
			}
			fmt.Printf("Search results exported to %s.\n", outputFile)
		} else {
			fmt.Println("Search Results:")
			for _, match := range matches {
				prettyJSON, _ := json.MarshalIndent(match, "", "  ")
				fmt.Println(string(prettyJSON))
			}
		}
	},
}

// Recursive function to search for matching objects
func searchObjects(data interface{}, key, value string) []map[string]interface{} {
	var results []map[string]interface{}

	switch obj := data.(type) {
	case map[string]interface{}:
		// Check if the current object contains the key-value pair
		if v, found := obj[key]; found && fmt.Sprint(v) == value {
			results = append(results, obj)
		}
		// Recursively search in nested objects
		for _, v := range obj {
			results = append(results, searchObjects(v, key, value)...)
		}
	case []interface{}:
		// Iterate through arrays and search each element
		for _, v := range obj {
			results = append(results, searchObjects(v, key, value)...)
		}
	}

	return results
}

func init() {
	// Add output flag to the search command
	searchCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Export results to a file")
	rootCmd.AddCommand(searchCmd)
}
