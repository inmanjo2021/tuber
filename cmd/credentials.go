package cmd

import (
	"io/ioutil"

	"github.com/spf13/viper"
)

func credentials() (creds []byte, err error) {
	viper.SetDefault("credentials-path", "/etc/tuber-credentials/credentials.json")
	credentialsPath := viper.GetString("credentials-path")
	creds, err = ioutil.ReadFile(credentialsPath)
	return
}
