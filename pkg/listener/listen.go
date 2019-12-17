package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"tuber/pkg/util"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// Listen it listens
func Listen(ctx context.Context) (in chan *util.RegistryEvent, out chan *util.RegistryEvent, err error) {
	in = make(chan *util.RegistryEvent, 1)
	out = make(chan *util.RegistryEvent, 1)

	err = startListener(ctx, in)
	if err != nil {
		return
	}
	go startAcker(out)
	return
}

func startListener(ctx context.Context, unprocessedEvents chan *util.RegistryEvent) error {
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
	go subscription.Receive(ctx,
		func(ctx context.Context, message *pubsub.Message) {
			obj := &util.RegistryEvent{Message: message}
			err := json.Unmarshal(message.Data, obj)
			if err != nil {
				fmt.Println("errors and stuff")
			} else {
				unprocessedEvents <- obj
			}
		})
	return err
}

func startAcker(processedEvents chan *util.RegistryEvent) {
	for event := range processedEvents {
		//event.Message.Ack()
		spew.Dump(event)
	}
}
