package errors

import (
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry initializes the Sentry error reporting client
func InitSentry(dsn string, errors <-chan error) (err error) {
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

	for err := range errors {
		alertSentry(err)
	}

	return nil
}

func alertSentry(err error) {
	sentry.CaptureException(err)
	sentry.Flush(time.Second * 5)
}
