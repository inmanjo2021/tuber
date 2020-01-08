package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"tuber/pkg/dataTemplate"
	"tuber/pkg/k8s"

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

// TODO: don't run this at the moment. need to add support for passing in interpolatables. it currently applies the local yamls raw.
func update(cmd *cobra.Command, args []string) {
	dir := ".tuber/"
	files, err := ioutil.ReadDir(".tuber")
	if err != nil {
		fmt.Println(err)
		return
	}

	var yamls []dataTemplate.Yaml

	for _, file := range files {
		name := file.Name()
		data, readErr := ioutil.ReadFile(dir + name)
		if readErr != nil {
			fmt.Println(readErr)
			return
		}
		content := string(data)
		yamls = append(yamls, dataTemplate.Yaml{Content: content, Filename: name})
	}

	lastIndex := len(yamls) - 1
	var buf bytes.Buffer

	for i, yaml := range yamls {
		_, err = io.WriteString(&buf, yaml.Content)

		if i < lastIndex {
			_, err = io.WriteString(&buf, "---\n")
		}
	}
	if err != nil {
		return
	}
	bytes := buf.Bytes()

	out, err := k8s.Apply(bytes, "tuber")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
