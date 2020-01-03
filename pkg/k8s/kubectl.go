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
func Get(kind string, name string, namespace string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "get", kind, name, "-o", "json", "-n", namespace)

	out, err = cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		err = fmt.Errorf(string(out))
	}

	return
}

// Apply applies Yaml vec
func Apply(yamls []util.Yaml, namespace string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "apply", "-n", namespace, "-f", "-")
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

// SetImage sets a digest for all deployments of a given namespaced container
func SetImage(namespace string, container string, digest string) (out []byte, err error) {
	cmd := exec.Command("kubectl", "-n", namespace, "set", "image", "deployments", container+"="+digest, "--all")

	out, err = cmd.CombinedOutput()

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
