package k8s

import "strings"

// CommandError wraps exitCode != 0 cases with kubectl and provides a richer interface
type CommandError struct {
	message string
}

func (e *CommandError) ResourceAlreadyExists() bool {
	return strings.Contains(e.message, "AlreadyExists")
}

func (e *CommandError) Error() string {
	return e.message
}
