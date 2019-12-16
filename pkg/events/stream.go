package events

import (
	"os"
	"tuber/pkg/release"
	"tuber/pkg/util"
)

func Stream(ch chan *util.RegistryEvent) {
	for event := range ch {
		event := filter(event)
		if event != nil {
			token := os.Getenv("GCLOUD_TOKEN")
			release.New(event.name, event.branch, token)
		}
	}
}
