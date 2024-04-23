package nats

import "errors"

var (
	ErrNATSNotConnected              = errors.New("nats not connected")
	ErrNATSServerHeadersNotSupported = errors.New("nats server headers not supported")
)

// Error godoc
type Error struct {
	Message string `json:"message"`
}
