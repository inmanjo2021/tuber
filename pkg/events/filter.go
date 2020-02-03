package events

import (
	"tuber/pkg/core"
	"tuber/pkg/listener"
)

func filter(e *listener.RegistryEvent) (event *core.TuberApp, err error) {
	apps, err := core.TuberApps()

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
