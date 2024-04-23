package v1

import (
	"github.com/ventive/go-mono-template/pkg/logger"
	"github.com/ventive/go-mono-template/pkg/nats"
	"github.com/ventive/go-mono-template/pkg/nats/middleware"
)

func (a *App) setupNats() error {
	log := logger.New(appID, "App.setupNats")

	log.Info("Connecting to NATS server")
	if err := a.nats.Connect(); err != nil {
		return err
	}

	if !a.nats.HeadersSupported() {
		log.Error("NATS server does not support headers", nats.ErrNATSServerHeadersNotSupported)

		return nats.ErrNATSServerHeadersNotSupported
	}

	if err := a.natsSubscribe(); err != nil {
		return err
	}

	return nil
}

func (a *App) natsSubscribe() error {
	middleware.LogSource = a.config.Logger.Source
	middleware.LogEnv = a.config.App.Env

	var err error

	queue := a.config.App.Queues.Subscribe.Queue
	if queue != "" {
		a.subscription, err = a.natsSubscribeTo(queue, a.subtractHandler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) natsSubscribeTo(queue string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	log := logger.New(appID, "App.natsSubscribeTo")

	log.Info("subscribing to " + queue)

	sub, err := a.nats.QueueSubscribe(queue, a.config.App.Nats.Name,
		middleware.UseMiddleware(handler, middleware.Log))
	if err != nil {
		log.Error("Error subscribing to "+queue, err)

		return nil, err
	}

	return sub, nil
}
