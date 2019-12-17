package events

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"tuber/pkg/release"
	"tuber/pkg/util"
)

// Stream streams a stream
func Stream(ch chan *util.RegistryEvent, token string) {
	for event := range ch {
		pendingRelease := filter(event)
		if pendingRelease == nil {
			event.Message.Ack()
			return
		}

		fmt.Println("Starting release for", pendingRelease.name, pendingRelease.branch)
		_, err := release.New(pendingRelease.name, pendingRelease.branch, token)

		if err != nil {
			spew.Dump(err)
		} else {
			event.Message.Ack()
		}
	}
}
