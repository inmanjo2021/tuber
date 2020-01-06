package events

import (
	"tuber/pkg/containers"
	"tuber/pkg/k8s"
	"tuber/pkg/pulp"
	"tuber/pkg/util"
)

// New create or update app in kubernetes
func publish(app *pulp.TuberApp, event *util.RegistryEvent, token string) (setImageOutput []byte, applyOutput []byte, err error) {
	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	containerName := event.ContainerName()

	yamls, err = updateImage(yamls, event)

	if err != nil {
		return
	}

	applyOutput, err = k8s.Apply(yamls, containerName)

	return
}
