package events

import (
	"sync"
	"tuber/pkg/listener"

	"go.uber.org/zap"
)

type streamer struct {
	token  string
	logger *zap.Logger
}

// NewStreamer creates a new Streamer struct
func NewStreamer(token string, logger *zap.Logger) *streamer {
	return &streamer{token, logger}
}

// Stream streams a stream
func (s *streamer) Stream(unprocessed <-chan *listener.RegistryEvent, processed chan<- *listener.RegistryEvent, chErr chan<- listener.FailedRelease) {
	defer close(processed)
	defer close(chErr)

	var wait = &sync.WaitGroup{}

	for event := range unprocessed {
		go func(event *listener.RegistryEvent) {
			wait.Add(1)
			defer wait.Done()

			var err error
			defer func() {
				if err != nil {
					chErr <- listener.FailedRelease{Err: err, Event: event}
				} else {
					processed <- event
				}
			}()

			pendingRelease, err := filter(event)

			if err != nil || pendingRelease == nil {
				return
			}

			var releaseLog = s.logger.With(
				zap.String("releaseName", pendingRelease.Name),
				zap.String("releaseBranch", pendingRelease.Tag))

			releaseLog.Info("release: starting")

			output, err := publish(pendingRelease, event.Digest, s.token)

			if err != nil {
				releaseLog.Warn(
					"release: error",
					zap.Error(err),
					zap.String("output", string(output)),
				)
			} else {
				releaseLog.Info("release: done")
			}
		}(event)
	}

	// Wait for all publish goroutines to be done.
	wait.Wait()
}
