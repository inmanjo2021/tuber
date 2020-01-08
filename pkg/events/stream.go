package events

import (
	"sync"
	"tuber/pkg/util"

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

type failedRelease struct {
	err   error
	event *util.RegistryEvent
}

// Err returns the error causing a failed release
func (f *failedRelease) Err() error { return f.err }

// Event returns the failed release event
func (f *failedRelease) Event() *util.RegistryEvent { return f.event }

// Stream streams a stream
func (s *streamer) Stream(unprocessed <-chan *util.RegistryEvent, processed chan<- *util.RegistryEvent, chErr chan<- util.FailedRelease) {
	defer close(processed)
	defer close(chErr)

	var wait = &sync.WaitGroup{}

	s.logger.Info("stream loop start")

	for event := range unprocessed {
		go func(event *util.RegistryEvent) {
			wait.Add(1)
			defer wait.Done()

			var err error
			defer func() {
				if err != nil {
					s.logger.Warn("failed release err", zap.Error(err))
					chErr <- &failedRelease{err: err, event: event}
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
