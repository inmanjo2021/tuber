package dataTemplate

import (
	"bytes"
	"text/template"
)

// Yaml is a yaml
type Yaml struct {
	Content  string
	Filename string
}

// Interpolate uses the standard template library
func (yaml Yaml) Interpolate(data map[string]string) (buf bytes.Buffer, err error) {
	tmpl, err := template.New("").Parse(yaml.Content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&buf, data)
	return
}

// NewInterpolated interpolates an existing yaml and returns a new one
func (yaml Yaml) NewInterpolated(data map[string]string) (interpolatedYaml Yaml, err error) {
	interpolated, err := yaml.Interpolate(data)
	interpolatedContent := interpolated.String()
	interpolatedYaml = Yaml{Content: interpolatedContent, Filename: yaml.Filename}
	return
}
