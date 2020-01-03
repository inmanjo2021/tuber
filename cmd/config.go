package cmd

import (
	"fmt"
	"tuber/pkg/k8s"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config [add]",
}

var configPatchCmd = &cobra.Command{
	Use:  "patch [appName] [key] [value]",
	Run:  configPatch,
	Args: cobra.ExactArgs(3),
}

func configPatch(cmd *cobra.Command, args []string) {
	appName := args[0]
	key := args[1]
	value := args[2]
	data := fmt.Sprintf(`{"data":{"%s":"%s"}}`, key, value)
	name := fmt.Sprintf("configmap/%s-config", appName)
	k8s.Patch(name, appName, data)
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPatchCmd)
}
