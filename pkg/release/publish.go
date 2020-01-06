package release

import (
	"tuber/pkg/containers"
	"tuber/pkg/k8s"
	"tuber/pkg/pulp"
	"tuber/pkg/util"
)

// New create or update app in kubernetes
func New(app *pulp.TuberApp, event *util.RegistryEvent, token string) (setImageOutput []byte, applyOutput []byte, err error) {
	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	containerName := event.ContainerName()

	applyOutput, err = k8s.Apply(yamls, app.Name)

	if err != nil {
		return
	}

	setImageOutput, err = k8s.SetImage(app.Name, containerName, event.Digest)

	return
}
