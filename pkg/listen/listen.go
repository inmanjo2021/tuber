package listen

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// RegistryEvent json deserialize target for pubsub
type RegistryEvent struct {
	Action string `json:"action"`
	Digest string `json:"digest"`
}

type Subscription struct {
	projectId     string
	subscription  string
	clientOptions []option.ClientOption
}

type SubscriptionOption func(*Subscription)

func WithCredentialsFile(credentials string) SubscriptionOption {
	return func(s *Subscription) {
		s.clientOptions = append(s.clientOptions, option.WithCredentialsFile(credentials))
	}

}

func NewSubscription(projectId string, subscription string, options ...SubscriptionOption) *Subscription {
	var s = &Subscription{
		projectId,
		subscription,
		[]option.ClientOption{},
	}
	for _, option := range options {
		option(s)
	}

	return s
}

// Listen it listens
func (s *Subscription) Listen(ctx context.Context, events chan *RegistryEvent) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()
	var client *pubsub.Client
	var err error

	client, err = pubsub.NewClient(ctx, "freshly-docker", option.WithCredentialsFile("./credentials.json"))

	if err != nil {
		client, err = pubsub.NewClient(ctx, "freshly-docker")
	}

	if err != nil {
		return err
	}

	subscription := client.Subscription(s.subscription)

	err = subscription.Receive(ctx,
		func(ctx context.Context, message *pubsub.Message) {
			var obj = new(RegistryEvent)
			err := json.Unmarshal(message.Data, obj)
			if err != nil {
				events <- obj
			} else {
				// log errors?
			}
		})
	return err
}
