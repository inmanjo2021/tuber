package listener

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"tuber/pkg/util"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)


type listener struct {
	logger *zap.Logger
	in chan *util.RegistryEvent
	out chan *util.RegistryEvent
}

// Listen it listens
func Listen(ctx context.Context, logger *zap.Logger) (<-chan *util.RegistryEvent, chan<- *util.RegistryEvent, error) {
	var l = &listener{
		logger: logger,
		in: make(chan *util.RegistryEvent, 1),
		out: make(chan *util.RegistryEvent, 1),
	}

	var err = l.startListener(ctx)

	if err != nil {
		return l.in, l.out, err
	}
	go l.startAcker()
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

	l.logger.Info("Listener: starting")
	go subscription.Receive(ctx,
		func(ctx context.Context, message *pubsub.Message) {
			obj := &util.RegistryEvent{Message: message}
			err := json.Unmarshal(message.Data, obj)
			if err != nil {
				l.logger.Warn("Could not unmarshal message")
			} else {
				l.in <- obj
			}
		})
	l.logger.Info("Listener: started")
	return err
}

func (l *listener) startAcker() {
	l.logger.Info("Acknowledge loop: starting")
	for event := range l.out {
		//event.Message.Ack()
		l.logger.With(zap.Reflect("event", event)).Info("Did not ack")
	}
	l.logger.Info("Acknowledge loop: stopped")
}
