package core

import (
	"bytes"
	"strings"
	"text/template"
	"tuber/pkg/k8s"
)

// ReleaseTubers combines and interpolates with tuber's conventions, and applies them
func ReleaseTubers(tubers []string, app *TuberApp, digest string) (output []byte, err error) {
	return ApplyInterpolated(tubers, app.Name, tuberData(app, digest))
}

func tuberData(app *TuberApp, digest string) (data map[string]string) {
	return map[string]string{
		"tuberImage": digest,
	}
}

// ApplyInterpolated combines and interpolates yamls with provided data, and applies them
func ApplyInterpolated(yamls []string, namespace string, data map[string]string) (output []byte, err error) {
	combined := strings.Join(yamls, "---\n")
	tmpl, err := template.New("").Parse(combined)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	interpolated := buf.Bytes()
	output, err = k8s.Apply(interpolated, namespace)
}
