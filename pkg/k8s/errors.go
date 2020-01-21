package k8s

import (
	"fmt"
	"regexp"
	"strings"
)

// CommandError wraps exitCode != 0 cases with kubectl and provides a richer interface
type CommandError struct {
	message string
}

func (e *CommandError) ResourceAlreadyExists() bool {
	return strings.Contains(e.message, "AlreadyExists")
}

func (e *CommandError) Error() string {
	re := regexp.MustCompile(`^.*: `)
	msg := re.ReplaceAllString(e.message, "")

	return fmt.Sprintf("k8s: %s", msg)
}
