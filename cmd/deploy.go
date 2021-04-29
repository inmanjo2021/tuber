package cmd

import (
	"context"
	"fmt"

	"github.com/freshly/tuber/pkg/containers"
	"github.com/freshly/tuber/pkg/core"
	"github.com/freshly/tuber/pkg/events"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "deploy [app]",
	Short:        "deploys the latest built image of an app",
	RunE:         deploy,
	PreRunE:      promptCurrentContext,
}

func deploy(cmd *cobra.Command, args []string) error {
	appName := args[0]
	logger, err := createLogger()
	if err != nil {
		return err
	}
	logger.With(zap.String("action", "deploy"))

	defer logger.Sync()

	logger.Debug("getting tuber source apps")
	apps, err := core.TuberSourceApps()
	if err != nil {
		return err
	}

	logger.Debug("getting credentials")
	creds, err := credentials()
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("attempting to find app named: %s", appName))
	app, err := apps.FindApp(appName)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("found app: %+v", app))

	logger.Debug("getting repository location")
	location := app.GetRepositoryLocation()

	logger.Debug("getting latest sha")
	sha, err := containers.GetLatestSHA(location, creds)
	if err != nil {
		return err
	}

	logger.Debug("getting cluster data")
	data, err := clusterData()
	if err != nil {
		return err
	}

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	logger.Debug("creating new processor")
	processor := events.NewProcessor(ctx, logger, creds, data, viper.GetBool("reviewapps-enabled"))
	digest := app.RepoHost + "/" + app.RepoPath + "@" + sha
	tag := app.ImageTag

	logger.Debug("processing message")
	processor.ProcessMessage(digest, tag)

	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
