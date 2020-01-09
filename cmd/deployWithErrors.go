package cmd

import (
	"fmt"
	"log"
	"tuber/pkg/containers"
	"tuber/pkg/core"
	"tuber/pkg/events"
	"tuber/pkg/gcloud"
	"tuber/pkg/util"

	"github.com/spf13/cobra"
)

var deployWithErrorsCmd = &cobra.Command{
	Use:   "deployWithErrors [appName]",
	Short: "Deploys an app and simulates an error",
	Run:   deployWithErrors,
	Args:  cobra.ExactArgs(1),
}

func deployWithErrors(cmd *cobra.Command, args []string) {
	logger, err := createLogger()

	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()

	apps, err := core.TuberApps()

	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials()
	if err != nil {
		log.Fatalln(err.Error())
	}

	app, err := apps.FindApp(args[0])
	if err != nil {
		log.Fatal(err)
	}

	location := app.GetRepositoryLocation()

	token, err := gcloud.GetAccessToken(creds)
	if err != nil {
		log.Fatal(err)
	}

	sha, err := containers.GetLatestSHA(location, token)
	if err != nil {
		log.Fatal(err)
	}

	streamer := events.NewStreamer(token, logger)

	errorChan := make(chan util.FailedRelease, 1)
	unprocessedEvents := make(chan *util.RegistryEvent, 1)
	processedEvents := make(chan *util.RegistryEvent, 1)
	errorReports := make(chan error, 1)
	go streamer.Stream(unprocessedEvents, processedEvents, errorChan, errorReports)

	ackable := emptyAckable{}
	deployEvent := util.RegistryEvent{
		Action:  "INSERT",
		Digest:  app.RepoHost + "/" + app.RepoPath + "@" + sha,
		Tag:     app.ImageTag,
		Message: ackable,
	}

	unprocessedEvents <- &deployEvent

	select {
	case msg := <-errorReports:
		fmt.Println("-------- error recieved ----------")
		fmt.Println(msg)
		close(errorReports)
	case <-errorChan:
		close(unprocessedEvents)
	case <-processedEvents:
		close(unprocessedEvents)
	}
}

func init() {
	rootCmd.AddCommand(deployWithErrorsCmd)
}
