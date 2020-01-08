package cmd

import (
	"log"
	"strings"
	"tuber/pkg/core"
	"tuber/pkg/k8s"

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

		err = k8s.CreateNamespace(appName)
		err = k8s.BindNamespace(appName)

		if err != nil && !strings.Contains(err.Error(), "AlreadyExists") {
			log.Fatal(err)
		}

		err = core.AddAppConfig(appName, repo, tag)

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
