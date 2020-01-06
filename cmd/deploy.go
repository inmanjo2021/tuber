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

type failedRelease struct {
	err   error
	event *util.RegistryEvent
}

// Err returns the error causing a failed release
func (f failedRelease) Err() error { return nil }

// Event returns the failed release event
func (f failedRelease) Event() *util.RegistryEvent { return f.event }

func deploy(cmd *cobra.Command, args []string) {
	logger := createLogger()
	defer logger.Sync()

	apps, err := pulp.TuberApps()

	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials()
	if err != nil {
		log.Fatal(err)
	}

	token, err := gcloud.GetAccessToken(creds)

	if err != nil {
		log.Fatal(err)
	}

	app, err := apps.FindApp(args[0])
	if err != nil {
		log.Fatal(err)
	}

	location := app.GetRepositoryLocation()

	sha, err := containers.GetLatestSHA(location, token)

	if err != nil {
		log.Fatal(err)
	}

	streamer := events.NewStreamer(token, logger)

	errorChan := make(chan util.FailedRelease, 1)
	unprocessedEvents := make(chan *util.RegistryEvent, 1)
	processedEvents := make(chan *util.RegistryEvent, 1)
	go streamer.Stream(unprocessedEvents, processedEvents, errorChan)

	ackable := emptyAckable{}
	deployEvent := util.RegistryEvent{
		Action:  "INSERT",
		Digest:  app.RepoHost + "/" + app.RepoPath + "@" + sha,
		Tag:     app.ImageTag,
		Message: ackable,
	}

	unprocessedEvents <- &deployEvent

	select {
	case <-errorChan:
		close(unprocessedEvents)
	case <-processedEvents:
		close(unprocessedEvents)
	}
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
