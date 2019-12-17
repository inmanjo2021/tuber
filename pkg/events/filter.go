package events

import (
	"fmt"

	"tuber/pkg/k8s"
	"tuber/pkg/util"
)

type pendingRelease struct {
	name   string
	branch string
}

func filter(e *util.RegistryEvent) (event pendingRelease, err error) {
	apps, err := k8s.TuberApps()
	var matchApp k8s.TuberApp
	found := false

	for _, app := range apps {
		if app.ImageTag == e.Tag {
			found = true
			matchApp = app
			break
		}
	}

	if !found {
		fmt.Println("Ignoring", e.Tag)
		return
	}

	event = pendingRelease{name: matchApp.Name, branch: matchApp.Tag}
	return
}
