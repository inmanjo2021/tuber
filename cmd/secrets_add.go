package cmd

import (
	"fmt"

	"tuber/pkg/k8s"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [local filepath] [namespace]",
	Short: "secrets add from file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		k8s.CreateFromFile(args[0], args[1])
	},
}

func init() {
	secretsCmd.AddCommand(addCmd)
}
