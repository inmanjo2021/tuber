package events

import (
	"bytes"
	"io"
	"tuber/pkg/containers"
	"tuber/pkg/dataTemplate"
	"tuber/pkg/k8s"
	"tuber/pkg/pulp"
	"tuber/pkg/util"
)

// New create or update app in kubernetes
func publish(app *pulp.TuberApp, event *util.RegistryEvent, token string) (output []byte, err error) {
	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	// TODO: make this smarter and separate
	data := map[string]string{"tuberImage": event.Digest}
	yamls, err = interpolate(yamls, data)

	if err != nil {
		return
	}

	yamlBytes, err := convert(yamls)

	if err != nil {
		return
	}

	output, err = k8s.Apply(yamlBytes, event.ContainerName())

	return
}

func convert(yamls []dataTemplate.Yaml) (out []byte, err error) {
	lastIndex := len(yamls) - 1
	var buf bytes.Buffer

	for i, yaml := range yamls {
		_, err = io.WriteString(&buf, yaml.Content)

		if i < lastIndex {
			_, err = io.WriteString(&buf, "---\n")
		}
	}
	out = buf.Bytes()
	return
}
