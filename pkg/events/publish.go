package events

import (
	"tuber/pkg/containers"
	"tuber/pkg/core"
)

func publish(app *core.TuberApp, digest string, token string, clusterData *core.ClusterData) (output []byte, err error) {
	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	output, err = core.ReleaseTubers(yamls, app, digest, clusterData)

	return
}
