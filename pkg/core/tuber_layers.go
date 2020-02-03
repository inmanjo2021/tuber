package core

import (
	"strings"
)

// ReleaseTubers combines and interpolates with tuber's conventions, and applies them
func ReleaseTubers(tubers []string, app *TuberApp, digest string, data *ClusterData) ([]byte, error) {
	return ApplyTemplate(app.Name, strings.Join(tubers, "---\n"), tuberData(app, digest, data))
}

// ClusterData is configurable, cluster-wide data available for yaml interpolation
type ClusterData struct {
	DefaultGateway  string
	DefaultHostname string
}

func tuberData(app *TuberApp, digest string, clusterData *ClusterData) (data map[string]string) {
	return map[string]string{
		"tuberImage":             digest,
		"clusterDefaultGateway":  clusterData.DefaultGateway,
		"clusterDefaultHostname": clusterData.DefaultHostname,
	}
}
