package cmd

import (
	"log"
	"tuber/pkg/core"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [app name] [docker repo] [deploy tag]",
	Short: "create new app in current cluster",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		appName := args[0]
		repo := args[1]
		tag := args[2]

		outputError, err := core.CreateTuberApp(appName, repo, tag)

		if err != nil {
			log.Fatal(err, string(outputError))
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
