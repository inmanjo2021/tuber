package listener

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

// FailedRelease can be created while streaming, consumed by the listener
type FailedRelease struct {
	Err   error
	Event *RegistryEvent
}
