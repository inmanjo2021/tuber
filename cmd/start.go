package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
	"tuber/pkg/events"
	"tuber/pkg/listener"
	"tuber/pkg/util"

	"context"
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
	var ch = make(chan *util.RegistryEvent, 20)
	token := os.Getenv("GCLOUD_TOKEN")

	go events.Stream(ch, token)
	go listener.Listen(ctx, ch)

	select {}
}
