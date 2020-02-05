package events

import (
	"fmt"
	"tuber/pkg/containers"
	"tuber/pkg/core"
)

func publish(app *core.TuberApp, digest string, creds []byte, clusterData *core.ClusterData) (err error) {
	prereleaseYamls, releaseYamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), creds)

	if err != nil {
		return
	}

	if len(prereleaseYamls) > 0 {
		err = core.RunPrerelease(prereleaseYamls, app, digest, clusterData)

		if err != nil {
			err = fmt.Errorf("prerelease error: %s", err.Error())
			return
		}
	}

	return core.ReleaseTubers(releaseYamls, app, digest, clusterData)
}
