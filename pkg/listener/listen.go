package listener

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
	"time"
	"tuber/pkg/util"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type listener struct {
	projectId    string
	subscription string

	in   chan *util.RegistryEvent
	out  chan *util.RegistryEvent
	wait *sync.WaitGroup

	logger       *zap.Logger
	recvSettings pubsub.ReceiveSettings
}

// ListenerOption provides optional settings to a listener constructor
type ListenerOption func(*listener)

// WithMaxAccept determines the maximum number of outstanding messages accepted
func WithMaxAccept(n int) ListenerOption {
	return func(l *listener) {
		l.recvSettings.MaxOutstandingMessages = n
	}
}

// WithMaxTimeout sets the maximum ack timeout extension
func WithMaxTimeout(d time.Duration) ListenerOption {
	return func(l *listener) {
		l.recvSettings.MaxExtension = d
	}
}

// NewListener creates a new PubSub listener
func NewListener(logger *zap.Logger, options ...ListenerOption) *listener {
	var l = &listener{
		projectId:    "freshly-docker",
		subscription: "freshly-docker-gcr-events",

		in:           make(chan *util.RegistryEvent, 1),
		out:          make(chan *util.RegistryEvent, 1),
		wait:         &sync.WaitGroup{},
		logger:       logger,
		recvSettings: pubsub.ReceiveSettings{},
	}

	for _, option := range options {
		option(l)
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
	var client *pubsub.Client
	var err error

	client, err = pubsub.NewClient(ctx, l.projectId, option.WithCredentialsFile("./credentials.json"))

	if err != nil {
		client, err = pubsub.NewClient(ctx, l.projectId)
	}

	if err != nil {
		return err
	}

	subscription := client.Subscription(l.subscription)
	subscription.ReceiveSettings = l.recvSettings

	go func(in chan<- *util.RegistryEvent, logger *zap.Logger) {
		// Register this goroutine in the waiter
		l.wait.Add(1)
		defer l.wait.Done()

		// Close the message channel before exiting to signal to downstream that we're done
		defer close(in)

		l.logger.Info("listener: starting")
		l.logger.Debug("listener: subscription options", zap.Reflect("options", subscription.ReceiveSettings))
		err = subscription.Receive(ctx,
			func(ctx context.Context, message *pubsub.Message) {
				obj := &util.RegistryEvent{Message: message}
				jsonErr := json.Unmarshal(message.Data, obj)

				if jsonErr != nil {
					l.logger.Warn("could not unmarshal message")
				} else {
					in <- obj
				}
			})

		if err != nil {
			l.logger.With(zap.Error(err)).Warn("listener: receiver error")
		}
		l.logger.Info("listener: shutting down")
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

	l.logger.Info("acknowledge loop: starting")

	for event := range l.out {
		//event.Message.Ack()
		l.logger.With(zap.Reflect("event", event)).Debug("did not ack")
	}
	l.logger.Info("acknowledge loop: stopped")
}

// Wait for the listener and acker goroutines to exit.
// If you use this method, you must ensure that you close
// the output channel when no more work is being processed
func (l *listener) Wait() {
	l.wait.Wait()
}
