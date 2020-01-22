package k8s

import (
	"errors"
	"strings"
)

var ErrResourceAlreadyExists = errors.New("k8s: resource already exists")
var ErrResourceNotFound = errors.New("k8s: resource not found")

func NewError(message string) error {
	if strings.Contains(message, "AlreadyExists") {
		return ErrResourceAlreadyExists
	}

	if strings.Contains(message, "doesn't have a resource") {
		return ErrResourceNotFound
	}

	return errors.New(message)
}
