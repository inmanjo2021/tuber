package core

import (
	"bytes"
	"html/template"
	"tuber/pkg/k8s"
)

// ApplyTemplate interpolates and applies a yaml to a given namespace
func ApplyTemplate(namespace string, templatestring string, params map[string]string) (out []byte, err error) {
	tpl, err := template.New("").Parse(templatestring)

	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, params)

	if err != nil {
		return
	}

	out, err = k8s.Apply(buf.Bytes(), namespace)

	return
}
