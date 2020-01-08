package k8s

import (
	"bytes"
	"strings"
	"text/template"
)

// ApplyYamls joins yamls, interpolates a single map of template data, and kubectl apply's
func ApplyYamls(yamls []string, data map[string]string, namespace string) (output []byte, err error) {
	processed, err := processYamls(yamls, data)
	if err != nil {
		return
	}
	output, err = Apply(processed, namespace)
	return
}

func processYamls(yamls []string, data map[string]string) (processed []byte, err error) {
	combined := strings.Join(yamls, "---\n")
	tmpl, err := template.New("").Parse(combined)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	processed = buf.Bytes()
	return
}
