package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tuber/pkg/containers"
	"tuber/pkg/pulp"
)

func init() {
	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch [image] [tag]",
	Short: "Fetch Tuber yaml files",
	RunE:  fetch,
	Args:  cobra.ExactArgs(2),
}

func fetch(cmd *cobra.Command, args []string) error {
	viper.BindEnv("gcloud-token", "GCLOUD_TOKEN")
	var token = viper.GetString("gcloud-token")

	apps, err := pulp.TuberApps()
	if err != nil {
		return err
	}
	app := apps.FindApp(args[0], args[1])
	if app == nil {
		return fmt.Errorf("not found %s:%s", args[0], args[1])
	}

	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err == nil {
		for i, yaml := range yamls {
			if i > 0 {
				fmt.Println("--")
			}
			fmt.Print(yaml)
		}
	}
	return err
}
