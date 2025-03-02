package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [file]",
	Short: "View JSON file as a collapsible tree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			fmt.Println("Invalid JSON:", err)
			return
		}

		app := tview.NewApplication()
		tree := tview.NewTreeView().SetRoot(createTree(jsonData, "root")).SetCurrentNode(nil)
		tree.SetBorder(true).SetTitle("JSON Viewer").SetTitleAlign(tview.AlignLeft)

		if err := app.SetRoot(tree, true).Run(); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func createTree(data interface{}, key string) *tview.TreeNode {
	node := tview.NewTreeNode(key)
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			node.AddChild(createTree(val, k))
		}
	case []interface{}:
		for i, val := range v {
			node.AddChild(createTree(val, fmt.Sprintf("[%d]", i)))
		}
	default:
		node.SetText(fmt.Sprintf("%s: %v", key, v))
	}
	return node
}
