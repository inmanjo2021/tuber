package listen

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// RegistryEvent json deserialize target for pubsub
type RegistryEvent struct {
	Action string `json:"action"`
	Digest string `json:"digest"`
	Tag string `json:"tag"`
}

type Subscription struct {
	projectId     string
	subscription  string
	clientOptions []option.ClientOption
}

// Listen it listens
func Listen(ctx context.Context, events chan *RegistryEvent) error {
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

	err = subscription.Receive(ctx,
		func(ctx context.Context, message *pubsub.Message) {
			var obj = new(RegistryEvent)
			err := json.Unmarshal(message.Data, obj)
			if err != nil {
				fmt.Println("errors and stuff")
			} else {
				events <- obj
			}
		})
	return err
}
