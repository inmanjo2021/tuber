package k8s

import (
	"errors"
	"strings"
)

var ResourceAlreadyExists = errors.New("k8s: resource already exists")
var ResourceNotFound = errors.New("k8s: resource not found")

func NewError(message string) error {
	if strings.Contains(message, "AlreadyExists") {
		return ResourceAlreadyExists
	}

	if strings.Contains(message, "doesn't have a resource") {
		return ResourceNotFound
	}

	return errors.New(message)
}
