package util

// Yaml is a yaml
type Yaml struct {
	Content  string
	Filename string
}

// RegistryEvent json deserialize target for pubsub
type RegistryEvent struct {
	Action string `json:"action"`
	Digest string `json:"digest"`
	Tag string `json:"tag"`
}
