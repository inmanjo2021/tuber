package events

import (
	"tuber/pkg/pulp"
	"tuber/pkg/util"
)

func filter(e *util.RegistryEvent) (event *pulp.TuberApp, err error) {
	apps, err := pulp.TuberApps()

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
