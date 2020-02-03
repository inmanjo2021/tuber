package cmd

import (
	"io/ioutil"
	"tuber/pkg/k8s"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func credentials() (creds []byte, err error) {
	viper.SetDefault("credentials-path", "/etc/tuber-credentials/credentials.json")
	credentialsPath := viper.GetString("credentials-path")
	creds, err = ioutil.ReadFile(credentialsPath)
	return
}

var credentialsCmd = &cobra.Command{
	Use:   "credentials [local filepath] [namespace]",
	Short: "add tuber secrets from file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = k8s.CreateTuberCredentials(args[0], args[1])
		return
	},
}

func init() {
	rootCmd.AddCommand(credentialsCmd)
}
