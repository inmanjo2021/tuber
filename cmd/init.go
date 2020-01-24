package cmd

import (
	"tuber/pkg/core"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [appName] [portNumber]",
	Short: "initialize a .tuber directory and relevant yamls",
	Long: `App name is the name of your app, which will be interpolated into the configuration files and used as the
		namespace, as well as other things. Port number is the port that your app will run on. The standard is port 80.`,
	Args: cobra.ExactArgs(2),
	RunE: initialize,
}

func initialize(cmd *cobra.Command, args []string) (err error) {
	appName := args[0]
	portNumber := args[1]

	return core.Init(appName, portNumber)
}

func init() {
	rootCmd.AddCommand(initCmd)
}
