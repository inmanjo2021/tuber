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
	SilenceUsage: true,
	Use:          "credentials [local filepath] [namespace]",
	Short:        "add tuber secrets from file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return k8s.CreateTuberCredentials(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(credentialsCmd)
}
