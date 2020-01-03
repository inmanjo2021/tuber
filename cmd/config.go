package cmd

import (
	"fmt"
	"log"
	"tuber/pkg/k8s"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config [patch || remove]",
}

var configPatchCmd = &cobra.Command{
	Use:  "patch [appName] [key] [value]",
	Run:  configPatch,
	Args: cobra.ExactArgs(3),
}

var configRemoveCmd = &cobra.Command{
	Use:  "remove [appName] [key]",
	Run:  configRemove,
	Args: cobra.ExactArgs(2),
}

func configPatch(cmd *cobra.Command, args []string) {
	appName := args[0]
	key := args[1]
	value := args[2]
	data := fmt.Sprintf(`{"data":{"%s":"%s"}}`, key, value)
	name := fmt.Sprintf("configmap/%s-config", appName)
	out, err := k8s.Patch(name, appName, data)
	println(string(out))
	if err != nil {
		log.Fatal(err)
	}
}

func configRemove(cmd *cobra.Command, args []string) {
	appName := args[0]
	key := args[1]
	data := fmt.Sprintf(`[{"op":"remove", "path":"/data/%s"}]`, key)
	name := fmt.Sprintf("configmap/%s-config", appName)
	out, err := k8s.Remove(name, appName, data)
	println(string(out))
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPatchCmd)
	configCmd.AddCommand(configRemoveCmd)
}
