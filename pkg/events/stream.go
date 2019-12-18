package events

import (
	"go.uber.org/zap"
	"tuber/pkg/release"
	"tuber/pkg/util"
)

type Streamer struct {
	token  string
	logger *zap.Logger
}

func NewStreamer(token string, logger *zap.Logger) *Streamer {
	return &Streamer{token, logger}
}

// Stream streams a stream
func (s *Streamer) Stream(chIn <-chan *util.RegistryEvent, chOut chan<- *util.RegistryEvent, chErr chan<- error) {
	for event := range chIn {
		pendingRelease, err := filter(event)

		if err != nil {
			chErr <- err
			continue
		}

		if pendingRelease == nil {
			chOut <- event
			continue
		}

		var releaseLog = s.logger.With(
			zap.String("releaseName", pendingRelease.name),
			zap.String("releaseBranch", pendingRelease.branch))

		go func() {
			releaseLog.Info("Release: starting")

			_, err = release.New(pendingRelease.name, pendingRelease.branch, s.token)
			if err != nil {
				releaseLog.Warn("Release: error", zap.Error(err))
				chErr <- err
			} else {
				releaseLog.Info("Release: done")
				chOut <- event
			}
		}()
	}

	close(chOut)
	close(chErr)
}
