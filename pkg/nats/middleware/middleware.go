package middleware

import (
	"github.com/nats-io/nats.go"
)

// Middleware function signature
type Middleware func(next nats.MsgHandler) nats.MsgHandler

// UseMiddleware when interact with NATS
func UseMiddleware(handler nats.MsgHandler, middleware ...Middleware) nats.MsgHandler {
	for _, m := range middleware {
		handler = m(handler)
	}

	return handler
}
