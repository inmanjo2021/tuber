package cmd

import (
	"fmt"
	"io/ioutil"
	"tuber/pkg/k8s"
	"tuber/pkg/util"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [app name]",
	Short: "apply local yams",
	Run:   update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update(cmd *cobra.Command, args []string) {
	dir := ".tuber/"
	files, err := ioutil.ReadDir(".tuber")
	if err != nil {
		fmt.Println(err)
		return
	}

	var yamls []util.Yaml

	for _, file := range files {
		name := file.Name()
		data, readErr := ioutil.ReadFile(dir + name)
		if readErr != nil {
			fmt.Println(readErr)
			return
		}
		content := string(data)
		yamls = append(yamls, util.Yaml{Content: content, Filename: name})
	}

	out, err := k8s.Apply(yamls, "tuber")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
