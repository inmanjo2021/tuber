package events

import (
	"go.uber.org/zap"
	"tuber/pkg/release"
	"tuber/pkg/util"
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
func (s *streamer) Stream(chIn <-chan *util.RegistryEvent, chOut chan<- *util.RegistryEvent, chErr chan<- error) {
	defer close(chOut)
	defer close(chErr)

	for event := range chIn {
		pending, err := filter(event)

		if err != nil {
			chErr <- err
			continue
		}

		if pending == nil {
			chOut <- event
			continue
		}

		var releaseLog = s.logger.With(
			zap.String("releaseName", pending.name),
			zap.String("releaseBranch", pending.branch))

		go func() {
			releaseLog.Info("Release: starting")

			_, err = release.New(pending.name, pending.branch, s.token)

			if err != nil {
				releaseLog.Warn("Release: error", zap.Error(err))
				chErr <- err
			} else {
				releaseLog.Info("Release: done")
				chOut <- event
			}
		}()
	}
}
