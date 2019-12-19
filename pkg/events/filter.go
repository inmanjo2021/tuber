package events

import (
	"tuber/pkg/k8s"
	"tuber/pkg/util"
)

func filter(e *util.RegistryEvent) (event *k8s.TuberApp, err error) {
	apps, err := k8s.TuberApps()

	if err != nil {
		return
	}

	for _, app := range apps {
		if app.ImageTag == e.Tag {
			return &app, nil
		}
	}
	return
}
