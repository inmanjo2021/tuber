package listen

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
	"tuber/pkg/yamldownloader"
)

// Listen it listens
func Listen() {
	ctx := context.Background()
	client, _ := pubsub.NewClient(ctx, "freshly-docker")
	subscription := client.Subscription("freshly-docker-gcr-events")

	err := subscription.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		yamldownloader.FindLayer()
	})

	if err != nil {
		log.Fatal(err) // Error will always be not nil. and not always an error.
	}
}
