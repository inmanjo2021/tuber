package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tuber",
		Short: "",
	}

	app string
)

// Execute executes
func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "app name")

	return rootCmd.Execute()
}
