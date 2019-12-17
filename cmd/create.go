package cmd

import (
	"github.com/spf13/cobra"
	"tuber/pkg/k8s"
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
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		repo := args[1]
		tag := args[2]

		k8s.AddAppConfig(appName, repo, tag)
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
