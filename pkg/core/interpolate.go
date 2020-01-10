package core

import (
	"strings"

	"github.com/spf13/viper"
)

// ReleaseTubers combines and interpolates with tuber's conventions, and applies them
func ReleaseTubers(tubers []string, app *TuberApp, digest string) ([]byte, error) {
	return ApplyTemplate(app.Name, strings.Join(tubers, "---\n"), tuberData(app, digest))
}

func tuberData(app *TuberApp, digest string) (data map[string]string) {
	return map[string]string{
		"tuberImage":            digest,
		"clusterDefaultGateway": viper.GetString("default-gateway"),
		"clusterHostname":       viper.GetString("hostname"),
	}
}
