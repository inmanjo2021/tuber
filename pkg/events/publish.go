package events

import (
	"tuber/pkg/containers"
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

	output, err = k8s.ApplyYamls(yamls, data, event.ContainerName())

	return
}
