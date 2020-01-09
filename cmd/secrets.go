package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "secrets has subcommands like add",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("secrets called")
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
}
