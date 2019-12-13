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

type callback func(*RegistryEvent, error)

// Listen it listens
func Listen(listener callback) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()
	var client *pubsub.Client
	var err error

	client, err = pubsub.NewClient(ctx, "freshly-docker", option.WithCredentialsFile("./credentials.json"))

	if err != nil {
		client, err = pubsub.NewClient(ctx, "freshly-docker")
	}

	if err != nil {
		log.Fatal(err)
	}

	subscription := client.Subscription("freshly-docker-gcr-events")

	err = subscription.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		var obj = new(RegistryEvent)
		err := json.Unmarshal(message.Data, obj)

		listener(obj, err)
	})

	return err
}
