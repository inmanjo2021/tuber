package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"tuber/pkg/pulp"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 3 {
			err = errors.New("Incorrect arguments")
			return
		}

		appName := args[0]
		repo := args[1]
		tag := args[2]

		err = pulp.AddAppConfig(appName, repo, tag)

		if err != nil {
			log.Fatal(err)
		}

		return nil
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
