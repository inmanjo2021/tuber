package k8s

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"tuber/pkg/util"
)

// write apply a string using kubectl
func write(bytes []byte) (out []byte, err error) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return
	}

	stdin.Write(bytes)
	stdin.Close()

	out, err = cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		log.Fatal(string(out))
	}

	if err != nil {
		return
	}

	return
}

// Get get a config
func Get(kind string, name string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "get", kind, name, "-o", "json")

	out, err = cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		println(out)
		println("SFDSFDSF")
		println(cmd.ProcessState.ExitCode())
		return nil, fmt.Errorf("%s", out)
	}

	if err != nil {
		return
	}

	return
}

// Apply applies Yaml vec
func Apply(yamls []util.Yaml) (out []byte, err error) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return
	}

	go func() {
		defer stdin.Close()
		lastIndex := len(yamls) - 1

		for i, yaml := range yamls {
			io.WriteString(stdin, yaml.Content)

			if i < lastIndex {
				io.WriteString(stdin, "---\n")
			}
		}
	}()

	out, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	return
}
