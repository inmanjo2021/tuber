package apply

import (
	"io"
	"os/exec"
	"tuber/pkg/util"
)

func Apply(yamls []util.Yaml) ([]byte, error) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	// cmd := exec.Command("cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
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

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return out, nil
}
