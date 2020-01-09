package errors

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func Stream(sentryEnabled bool, sentryDSN string, errors <-chan error, logger *zap.Logger) {
	if sentryEnabled {
		err := InitSentry(sentryDSN)

		if err != nil {
			logger.Warn("failed to init sentry", zap.Error(err))
			sentryEnabled = false
		}
	}

	for err := range errors {
		logger.Warn("nonspecific error", zap.Error(err))

		if sentryEnabled {
			alertSentry(err)
		}
	}
}

// InitSentry initializes the Sentry error reporting client
func InitSentry(dsn string) (err error) {
	err = sentry.Init(
		sentry.ClientOptions{
			Dsn:              dsn,
			AttachStacktrace: true,
		},
	)

	if err != nil {
		return
	}

	defer sentry.Recover()

	return nil
}

func alertSentry(err error) {
	sentry.CaptureException(err)
	sentry.Flush(time.Second * 5)
}
