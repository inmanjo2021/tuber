package cmd

import (
	"fmt"
	"tuber/pkg/containers"
	"tuber/pkg/gcloud"
	"tuber/pkg/pulp"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch [appName]",
	Short: "Fetch Tuber yaml files",
	RunE:  fetch,
	Args:  cobra.ExactArgs(1),
}

func fetch(cmd *cobra.Command, args []string) error {
	token, err := gcloud.GetAccessToken()

	if err != nil {
		return err
	}

	apps, err := pulp.TuberApps()

	if err != nil {
		return err
	}

	app := apps.FindApp(args[0])

	if app == nil {
		return fmt.Errorf("app not found %s", args[0])
	}

	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err == nil {
		for i, yaml := range yamls {
			if i > 0 {
				fmt.Println("--")
			}
			fmt.Printf("%s\n", yaml)
		}
	}
	return err
}
