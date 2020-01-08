package core

import (
	"bytes"
	"strings"
	"text/template"
)

// ProcessYamls combines and interpolates with tuber's conventions
func ProcessYamls(yamls []string, app *TuberApp, digest string) (processed []byte, err error) {
	combined := strings.Join(yamls, "---\n")
	tmpl, err := template.New("").Parse(combined)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, interpolatableData(app, digest))
	processed = buf.Bytes()
	return
}

func interpolatableData(app *TuberApp, digest string) (data map[string]string) {
	return map[string]string{
		"tuberImage": digest,
	}
}
