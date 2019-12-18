package listener

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"sync"
	"tuber/pkg/util"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type listener struct {
	logger *zap.Logger
	in     chan *util.RegistryEvent
	out    chan *util.RegistryEvent
	wait   *sync.WaitGroup
}

// NewListener creates a new PubSub listener
func NewListener(logger *zap.Logger) *listener {
	var l = &listener{
		logger: logger,
		in:     make(chan *util.RegistryEvent, 1),
		out:    make(chan *util.RegistryEvent, 1),
		wait:   &sync.WaitGroup{},
	}
	return l
}

// Listen for incoming pubsub requests
func (l *listener) Listen(ctx context.Context) (<-chan *util.RegistryEvent, chan<- *util.RegistryEvent, error) {
	go l.startAcker(ctx)

	var err = l.startListener(ctx)
	return l.in, l.out, err
}

func (l *listener) startListener(ctx context.Context) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var client *pubsub.Client
	var err error

	client, err = pubsub.NewClient(ctx, "freshly-docker", option.WithCredentialsFile("./credentials.json"))

	if err != nil {
		client, err = pubsub.NewClient(ctx, "freshly-docker")
	}

	if err != nil {
		return err
	}

	subscription := client.Subscription("freshly-docker-gcr-events")
	cfg, err := subscription.Config(ctx)

	if err != nil {
		return err
	}

	// Set a fixed deadline and small outstanding message count for testing
	// TODO: Remove or make configurable
	subscription.ReceiveSettings.MaxExtension = cfg.AckDeadline
	subscription.ReceiveSettings.MaxOutstandingMessages = 5

	go func(in chan<- *util.RegistryEvent, logger *zap.Logger) {
		// Register this goroutine in the waiter
		l.wait.Add(1)
		defer l.wait.Done()

		// Close the message channel before exiting to signal to downstream that we're done
		defer close(in)

		l.logger.Info("Listener: starting")
		err = subscription.Receive(ctx,
			func(ctx context.Context, message *pubsub.Message) {
				obj := &util.RegistryEvent{Message: message}
				jsonErr := json.Unmarshal(message.Data, obj)

				if jsonErr != nil {
					l.logger.Warn("Could not unmarshal message")
				} else {
					in <- obj
				}
			})

		if err != nil {
			l.logger.With(zap.Error(err)).Warn("Listener: Receiver error")
		}
		l.logger.Info("Listener: shutting down")
	}(l.in, l.logger)

	return err
}

func (l *listener) startAcker(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	// Register this goroutine in the waiter
	l.wait.Add(1)
	defer l.wait.Done()

	l.logger.Info("Acknowledge loop: starting")

	for event := range l.out {
		//event.Message.Ack()
		l.logger.With(zap.Reflect("event", event)).Info("Did not ack")
	}
	l.logger.Info("Acknowledge loop: stopped")
}

// Wait for the listener and acker goroutines to exit.
// If you use this method, you must ensure that you close
// the output channel when no more work is being processed
func (l *listener) Wait() {
	l.wait.Wait()
}
