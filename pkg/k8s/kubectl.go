package k8s

import (
	"os/exec"
)

func runKubectl(cmd *exec.Cmd) (out []byte, err error) {
	out, err = cmd.CombinedOutput()

	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		err = newK8sError(out, err)
	}
	return
}

func kubectl(args ...string) ([]byte, error) {
	return runKubectl(exec.Command("kubectl", args...))
}

func pipeToKubectl(data []byte, args ...string) (out []byte, err error) {
	cmd := exec.Command("kubectl", args...)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return
	}

	_, err = stdin.Write(data)
	if err != nil {
		return
	}

	err = stdin.Close()
	if err != nil {
		return
	}

	return runKubectl(cmd)
}

// Apply `kubectl apply` data to a given namespace. Specify output or any other flags as args.
// Uses a stdin pipe to include the content of the data slice
func Apply(data []byte, namespace string, args ...string) ([]byte, error) {
	apply := []string{"apply", "-n", namespace, "-f", "-"}
	return pipeToKubectl(data, append(apply, args...)...)
}

// Get `kubectl get` a resource. Specify output or any other flags as args
func Get(kind string, name string, namespace string, args ...string) ([]byte, error) {
	get := []string{"get", kind, name, "-n", namespace}
	return kubectl(append(get, args...)...)
}

// Delete `kubectl delete` a resource. Specify output or any other flags as args
func Delete(kind string, name string, namespace string, args ...string) ([]byte, error) {
	deleteArgs := []string{"delete", kind, name, "-n", namespace}
	return kubectl(append(deleteArgs, args...)...)
}

// Create `kubectl create` a resource.
// Some resources take multiple args (like secrets), so both the resource type and any flags are the variadic
func Create(namespace string, resourceAndArgs ...string) ([]byte, error) {
	create := []string{"create", "-n", namespace}
	return kubectl(append(create, resourceAndArgs...)...)
}
