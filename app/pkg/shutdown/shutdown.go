package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Operation func(ctx context.Context) error

// GracefulShutdown waits for termination syscalls and doing clean up operations after received it
func GracefulShutdown(ctx context.Context, operations map[string]Operation, opts ...Option) <-chan struct{} {
	s := settings{
		timeout: 5 * time.Second,
		logger:  defaultLogger,
	}
	for _, opt := range opts {
		opt(&s)
	}

	wait := make(chan struct{})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		sig := <-sigChan

		s.logger.Infof("Received signal %s, shutting down", sig)

		// set timeout for the operations to be done to prevent system hang
		timeoutFunc := time.AfterFunc(s.timeout, func() {
			s.logger.Errorf("Timeout %s has been elapsed, force exit", s.timeout)
			os.Exit(1)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for k, v := range operations {
			wg.Add(1)

			go func(name string, operation Operation) {
				defer wg.Done()

				s.logger.Infof("Cleaning up: %s", name)
				if err := operation(ctx); err != nil {
					s.logger.Errorf("%s: clean up failed: %s", name, err)
					return
				}

				s.logger.Infof("%s was shutdown gracefully", name)
			}(k, v)
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
