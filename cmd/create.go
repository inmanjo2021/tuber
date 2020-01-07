package cmd

import (
	"log"
	"strings"
	"tuber/pkg/k8s"
	"tuber/pkg/pulp"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [app name] [docker repo] [deploy tag]",
	Short: "create new app in current cluster",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(3),
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

		err = pulp.AddAppConfig(appName, repo, tag)

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
