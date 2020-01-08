package events

import (
	"tuber/pkg/containers"
	"tuber/pkg/core"
	"tuber/pkg/k8s"
)

func publish(app *core.TuberApp, digest string, token string) (output []byte, err error) {
	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	processedYamls, err := core.ProcessYamls(yamls, app, digest)

	if err != nil {
		return
	}

	output, err = k8s.Apply(processedYamls, app.Name)

	return
}
