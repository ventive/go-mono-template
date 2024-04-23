package v1

import (
	"context"

	"github.com/ventive/go-mono-template/pkg/logger"
	"github.com/ventive/go-mono-template/pkg/nats"
)

type App struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	config       config
	nats         nats.Client
	subscription *nats.Subscription
}

func New(parentCtx context.Context, cfg config) (*App, error) {
	log := logger.New(appID, "New")
	ctx, cancelFunc := context.WithCancel(parentCtx)
	app := &App{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		config:     cfg,
	}
	log.Info("Setting up the app dependencies...")

	log.Info("Setting up NATS client")
	app.nats = nats.NewClient(nats.Config{
		URL:  cfg.App.Nats.URL,
		Name: cfg.App.Nats.Name,
		User: cfg.App.Nats.User,
		Pass: cfg.App.Nats.Pass,
		TLS: nats.TLSConfig{
			Enabled: cfg.App.Nats.TLS.Enabled,
			Cert:    cfg.App.Nats.TLS.Cert,
			Key:     cfg.App.Nats.TLS.Key,
			CA:      cfg.App.Nats.TLS.CA,
		},
	})
	return app, nil
}

// Start the application.
func (a *App) Start() error {
	defer a.cleanup()
	log := logger.New(appID, "App.Start")
	log.Debug("Starting up")

	log.Debug("setting up NATS")
	if err := a.setupNats(); err != nil {
		return err
	}

	<-a.ctx.Done()

	return a.ctx.Err()
}

// Stop application by calling the app's context cancelFunc.
func (a *App) Stop() {
	a.cancelFunc()
}

func (a *App) cleanup() {
	log := logger.New(appID, "App.cleanup")
	log.Debug("Cleanup started...")

	log.Info("Subscription: shutting down subscription")
	if err := a.subscription.Unsubscribe(); err != nil {
		log.ErrorWithExtra("Could not unsubscribe", map[string]interface{}{
			"service": "process",
			"action":  "nats:unsubscribe",
			"queue":   a.config.App.Queues.Subscribe.Queue,
		}, err)
	}
	log.Info("Subscription: draining subscription")
	if err := a.subscription.Drain(); err != nil {
		log.ErrorWithExtra("Failed to drain subscribtion", map[string]interface{}{
			"service": "process",
			"action":  "nats:sub:drain",
			"queue":   a.config.App.Queues.Subscribe.Queue,
		}, err)
	}

	log.Info("Closing NATS connection")
	a.nats.Close()

	log.Debug("Cleanup finished")
}
