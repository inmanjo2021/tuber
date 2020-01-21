package k8s

import (
	"os/exec"
)

// CreateNamespace create a new namespace in kubernetes
func CreateNamespace(namespace string) (err error) {
	cmd := exec.Command("kubectl", "create", "namespace", namespace)

	out, err := cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		err = &CommandError{message: string(out)}
	}

	return
}
