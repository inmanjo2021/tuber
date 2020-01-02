package util

import "strings"

// Yaml is a yaml
type Yaml struct {
	Content  string
	Filename string
}

// RegistryEvent json deserialize target for pubsub
type RegistryEvent struct {
	Action  string `json:"action"`
	Digest  string `json:"digest"`
	Tag     string `json:"tag"`
	Message ackable
}

type ackable interface {
	Ack()
	Nack()
}

// ContainerName extracts the container name from the tag
func (r *RegistryEvent) ContainerName() (name string) {
	tagSplit := strings.Split(r.Tag, "/")
	containerTag := tagSplit[len(tagSplit)-1]
	containerTagSplit := strings.Split(containerTag, ":")
	name = containerTagSplit[0]
	return
}
