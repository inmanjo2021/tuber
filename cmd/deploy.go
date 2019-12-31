package cmd

import (
	"log"
	"tuber/pkg/containers"
	"tuber/pkg/events"
	"tuber/pkg/gcloud"
	"tuber/pkg/pulp"
	"tuber/pkg/util"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [appName]",
	Short: "Deploys an app",
	Run:   deploy,
	Args:  cobra.ExactArgs(1),
}

type emptyAckable struct{}

func (emptyAckable) Ack()  {}
func (emptyAckable) Nack() {}

func deploy(cmd *cobra.Command, args []string) {
	logger := createLogger()
	defer logger.Sync()
	var errorChan = createErrorChannel(logger)

	apps, err := pulp.TuberApps()

	if err != nil {
		log.Fatal(err)
	}

	token, err := gcloud.GetAccessToken()

	if err != nil {
		log.Fatal(err)
	}

	app := apps.FindApp(args[0])
	location := app.GetRepositoryLocation()

	sha, err := containers.GetLatestSHA(location, token)

	if err != nil {
		log.Fatal(err)
	}

	streamer := events.NewStreamer(token, logger)

	unprocessedEvents := make(chan *util.RegistryEvent, 1)
	processedEvents := make(chan *util.RegistryEvent, 1)
	go streamer.Stream(unprocessedEvents, processedEvents, errorChan)

	ackable := emptyAckable{}
	deployEvent := util.RegistryEvent{
		Action:  "INSERT",
		Digest:  "gcr.io/freshly-docker/tuber@" + sha,
		Tag:     "gcr.io/freshly-docker/tuber:master",
		Message: ackable,
	}

	unprocessedEvents <- &deployEvent

	for range processedEvents {
		close(unprocessedEvents)
	}
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
