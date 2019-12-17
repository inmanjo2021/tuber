package cmd

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
	"tuber/pkg/events"
	"tuber/pkg/listener"
)

func init() {
	rootCmd.AddCommand(startCmd)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tuber",
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {
	var ctx = context.Background()

	unprocessedEvents, processedEvents, err := listener.Listen(ctx)
	if err != nil {
		panic(err)
	}

	var errorChan = make(chan error, 1)

	streamer := &events.Streamer{Token: os.Getenv("GCLOUD_TOKEN")}
	go func() {
		for error := range errorChan {
			log.Fatal(error)
		}
	}()
	go streamer.Stream(unprocessedEvents, processedEvents, errorChan)

	select {}
}
