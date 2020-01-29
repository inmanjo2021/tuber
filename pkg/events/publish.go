package events

import (
	"tuber/pkg/containers"
	"tuber/pkg/core"
)

func publish(app *core.TuberApp, digest string, token string) (output []byte, err error) {
	prereleaseYamls, releaseYamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	output, err = core.RunPrerelease(prereleaseYamls, app, digest)

	if err != nil {
		return
	}

	output, err = core.ReleaseTubers(releaseYamls, app, digest)

	return
}
