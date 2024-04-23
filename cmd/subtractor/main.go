package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	app "github.com/ventive/go-mono-template/services/subtractor/app/v1"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		errCh <- app.Run(ctx)
		close(errCh)
	}()

	select {
	case <-c:
	case <-errCh:
	}

	cancelFunc()
	wg.Wait()

}
