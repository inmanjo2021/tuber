package cmd

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"context"
	"tuber/pkg/listen"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tuber",
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {
	godotenv.Load()
	var ctx = context.Background()

	var ch = make(chan *listen.RegistryEvent, 20)
	var s = listen.NewSubscription("freshly-docker",
		"freshly-docker-gcr-events",
		listen.WithCredentialsFile("./credentials.json"))
	s.Listen(ctx, ch)

	for event := range ch {
		spew.Dump(event)
	}
}
