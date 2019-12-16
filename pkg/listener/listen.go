package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tuber/pkg/util"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Subscription struct {
	projectId     string
	subscription  string
	clientOptions []option.ClientOption
}

// Listen it listens
func Listen(ctx context.Context, events chan *util.RegistryEvent) error {
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

	fmt.Println("Listening...")
	err = subscription.Receive(ctx,
		func(ctx context.Context, message *pubsub.Message) {
			obj := &util.RegistryEvent { Message: message }
			err := json.Unmarshal(message.Data, obj)
			if err != nil {
				fmt.Println("errors and stuff")
			} else {
				events <- obj
			}
		})
	return err
}
