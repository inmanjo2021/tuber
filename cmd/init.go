package cmd

import (
	"tuber/pkg/core"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [appName] [routePrefix]",
	Short: "initialize a .tuber directory and relevant yamls",
	Long: `App name is the name of your app, which will be interpolated into the configuration files and used as the
		namespace, as well as other things.`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(2),
	RunE:         initialize,
}

func initialize(cmd *cobra.Command, args []string) (err error) {
	appName := args[0]
	routePrefix := args[1]

	return core.InitTuberApp(appName, routePrefix)
}

func init() {
	rootCmd.AddCommand(initCmd)
}
