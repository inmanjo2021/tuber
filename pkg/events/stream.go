package events

import (
	"fmt"
	"tuber/pkg/release"
	"tuber/pkg/util"
)

type Streamer struct {
	Token string
}

// Stream streams a stream
func (s *Streamer) Stream(chIn <-chan *util.RegistryEvent, chOut chan<- *util.RegistryEvent, chErr chan<- error) {
	for event := range chIn {
		pendingRelease, err := filter(event)

		if err != nil {
			chErr <- err
		}

		if pendingRelease == nil {
			chOut <- event
			continue
		}

		fmt.Println("Starting release for", pendingRelease.name, pendingRelease.branch)
		go func() {
			_, err = release.New(pendingRelease.name, pendingRelease.branch, s.Token)
			if err != nil {
				chErr <- err
			} else {
				chOut <- event
			}
		}()
	}
}
