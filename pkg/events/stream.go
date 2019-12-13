package events

import (
	"os"
	"tuber/pkg/release"
	"tuber/pkg/util"
)

func Stream(ch chan *util.RegistryEvent) {
	for event := range ch {
		qualified, e := filter(event)
		if qualified {
			token := os.Getenv("GCLOUD_TOKEN")
			release.New(e.name, e.branch, token)
		}
	}
}
