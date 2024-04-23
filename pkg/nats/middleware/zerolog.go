package middleware

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// LogSource godoc
	LogSource = "default-log-source"
	// LogEnv godoc
	LogEnv = ""
)

// Log reads the headers from NATS message and log them using zerolog pkg
func Log(next nats.MsgHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		zerologHeaderParams := zerolog.Dict()
		for k, v := range msg.Header {
			zerologHeaderParams = zerologHeaderParams.Interface(string(k), v)
		}

		log.Debug().Str("name", LogSource).
			Str("env", LogEnv).
			Dict("headers", zerologHeaderParams).Msg(fmt.Sprintf("Processing message from=%s", msg.Subject))

		next(msg)
	}
}
