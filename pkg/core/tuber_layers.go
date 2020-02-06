package core

import (
	"bytes"
	"strings"
	"text/template"
)

// ReleaseTubers combines and interpolates with tuber's conventions, and applies them
func ReleaseTubers(tubers []string, app *TuberApp, digest string, data *ClusterData) error {
	return ApplyTemplate(app.Name, strings.Join(tubers, "---\n"), tuberData(digest, data))
}

// ClusterData is configurable, cluster-wide data available for yaml interpolation
type ClusterData struct {
	DefaultGateway string
	DefaultHost    string
}

func tuberData(digest string, clusterData *ClusterData) (data map[string]string) {
	return map[string]string{
		"tuberImage":            digest,
		"clusterDefaultGateway": clusterData.DefaultGateway,
		"clusterDefaultHost":    clusterData.DefaultHost,
	}
}

func interpolateTuber(templateString string, data map[string]string) (interpolated []byte, err error) {
	tpl, err := template.New("").Parse(templateString)

	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)

	if err != nil {
		return
	}

	interpolated = buf.Bytes()
	return
}
