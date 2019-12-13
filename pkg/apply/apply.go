package apply

import (
	"io"
	"os/exec"
	"tuber/pkg/util"
)

// Write apply a string using kubectl
func Write(bytes []byte) (out []byte, err error) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	stdin, err := cmd.StdinPipe()
	defer stdin.Close()

	if err != nil {
		return
	}

	stdin.Write(bytes)
	out, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	return
}

// Apply applies Yaml vec
func Apply(yamls []util.Yaml) (out []byte, err error) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	// cmd := exec.Command("cat")
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
