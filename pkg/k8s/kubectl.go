package k8s

import (
	"fmt"
	"os/exec"
)

// Apply `kubectl apply` data to a given namespace
func Apply(bytes []byte, namespace string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "apply", "-n", namespace, "-f", "-")
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
		return nil, fmt.Errorf(string(out))
	}

	if cmd.ProcessState.ExitCode() != 0 {
		err = NewError(string(out))
	}

	return
}

// Get get a config
func Get(kind string, name string, namespace string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "get", kind, name, "-o", "json", "-n", namespace)

	out, err = cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		err = NewError(string(out))
	}

	return
}

// Patch patches data for a given resource and namespace
func Patch(name string, namespace string, data string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "patch", name, "-n", namespace, "--type", "merge", "-p", data)

	out, err = cmd.CombinedOutput()

	return
}

// Remove expects a remove operation with a path
func Remove(name string, namespace string, data string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "patch", name, "-n", namespace, "--type=json", "-p", data)

	out, err = cmd.CombinedOutput()

	return
}
