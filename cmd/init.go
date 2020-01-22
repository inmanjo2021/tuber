package cmd

import (
	"tuber/pkg/core"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [appName]",
	Short: "initialize a .tuber directory and relevant yamls",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE:  initialize,
}

func initialize(cmd *cobra.Command, args []string) (err error) {
	appName := args[0]

	if err = core.CreateTuberDirectory(); err != nil {
		return
	}

	if err = core.CreateDeploymentYAML(appName); err != nil {
		return
	}

	return
}

func init() {
	rootCmd.AddCommand(initCmd)
}
