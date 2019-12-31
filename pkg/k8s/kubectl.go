package k8s

import (
	"fmt"
	"io"
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

	_, err = stdin.Write(bytes)
	if err != nil {
		return
	}

	err = stdin.Close()
	if err != nil {
		return
	}

	out, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	if cmd.ProcessState.ExitCode() != 0 {
		err = fmt.Errorf(string(out))
	}

	return
}

// Get get a config
func Get(kind string, name string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "get", kind, name, "-o", "json")

	out, err = cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		err = fmt.Errorf(string(out))
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

	lastIndex := len(yamls) - 1

	for i, yaml := range yamls {
		io.WriteString(stdin, yaml.Content)

		if i < lastIndex {
			io.WriteString(stdin, "---\n")
		}
	}

	stdin.Close()
	out, err = cmd.CombinedOutput()

	return
}
