package v1

import (
	"context"
	"fmt"
	"sync"

	"github.com/ventive/go-mono-template/pkg/cli"
	"github.com/ventive/go-mono-template/pkg/logger"
	"github.com/ventive/go-mono-template/pkg/version"
)

const appID = "adder"

var configFile string

func Run(ctx context.Context) error {
	cli.Init(appID, "adder service")
	_ = cli.AddCommand("version", "Get the application version and Git commit SHA", logVersionDetails)
	_ = cli.AddCommand("start", "Start the service", start)
	cli.AssignStringFlag(&configFile, "config", "", "config file (default is ./.config.yaml)")

	return cli.Run(ctx)
}

func start(parentCtx context.Context) {
	log := logger.New(appID, "start")

	cfg, err := newConfig()
	if err != nil {
		log.Error("Unable to initialize config", err)
	}

	logger.Init(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})

	ctx, cancelFunc := context.WithCancel(parentCtx)
	app, err := New(ctx, cfg)
	if err != nil {
		log.Error("Unable to initialize adder", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer cancelFunc()
		log.Info("Starting adder")
		log.Error("adder stopped", app.Start())
	}()

	<-ctx.Done()

	app.Stop()

	log.Debug("Context is done. Waiting for WaitGroup to be done")
	wg.Wait()
}

func logVersionDetails(_ context.Context) {
	log := logger.New(appID, "logVersionDetails")
	log.Info(fmt.Sprintf("AppVersion=%s, GitCommit=%s", version.AppVersion, version.GitCommit))
}
