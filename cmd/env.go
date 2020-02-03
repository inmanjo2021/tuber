package cmd

import (
	"fmt"
	"log"
	"tuber/pkg/k8s"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use: "env [set || unset || file]",
}

var envSetCmd = &cobra.Command{
	Use:  "set [appName] [key] [value]",
	Run:  envSet,
	Args: cobra.ExactArgs(3),
}

var envUnsetCmd = &cobra.Command{
	Use:  "unset [appName] [key]",
	Run:  envUnset,
	Args: cobra.ExactArgs(2),
}

var fileCmd = &cobra.Command{
	Use:   "file [app] [local filepath]",
	Short: "batch env set",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return k8s.CreateEnvFromFile(args[0], args[1])
	},
}

func envSet(cmd *cobra.Command, args []string) {
	appName := args[0]
	key := args[1]
	value := args[2]
	mapName := fmt.Sprintf("%s-env", appName)
	err := k8s.PatchSecret(mapName, appName, key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func envUnset(cmd *cobra.Command, args []string) {
	appName := args[0]
	key := args[1]
	mapName := fmt.Sprintf("%s-env", appName)
	err := k8s.RemoveSecretEntry(mapName, appName, key)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envUnsetCmd)
	envCmd.AddCommand(fileCmd)
}
