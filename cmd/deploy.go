/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"tuber/pkg/containers"
	"tuber/pkg/gcloud"
	"tuber/pkg/pulp"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [appName]",
	Short: "Deploys an app",
	Run:   deploy,
	Args:  cobra.ExactArgs(1),
}

func deploy(cmd *cobra.Command, args []string) {
	apps, err := pulp.TuberApps()

	if err != nil {
		log.Fatal(err)
	}

	token, err := gcloud.GetAccessToken()

	if err != nil {
		log.Fatal(err)
	}

	app := apps.FindApp(args[0])
	location := app.GetRepositoryLocation()
	registry := containers.NewRegistry(location.Host, token)
	repository, err := registry.GetRepository(location.Path)

	if err != nil {
		log.Fatal(err)
	}

	manifest, err := repository.GetManifest(app.Tag)

	if err != nil {
		log.Fatal(err)
	}

	// repo location
	//   registry
	//     repository
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
