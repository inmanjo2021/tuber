package core

import (
	"strings"
	"tuber/pkg/k8s"
)

// ReleaseTubers combines and interpolates with tuber's conventions, and applies them
func ReleaseTubers(tubers []string, app *TuberApp, digest string) ([]byte, error) {
	return k8s.ApplyTemplate(app.Name, strings.Join(tubers, "---\n"), tuberData(app, digest))
}

func tuberData(app *TuberApp, digest string) (data map[string]string) {
	return map[string]string{
		"tuberImage": digest,
	}
}
