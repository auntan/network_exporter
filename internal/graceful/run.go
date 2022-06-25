package graceful

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const shutdownTimeout = time.Second * 5

func Run(ctx context.Context, runnables map[string]func(context.Context) error) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(runnables))

	type state struct {
		err  error
		done bool
	}
	states := make(map[string]*state, len(runnables))
	var xStates sync.Mutex

	for name, runnable := range runnables {
		go func(name string, f func(context.Context) error) {
			defer wg.Done()

			err := f(cancelCtx)

			xStates.Lock()
			states[name] = &state{
				err:  err,
				done: true,
			}
			xStates.Unlock()

			cancel()

			//zap.S().Debugf("%s done", name)
		}(name, runnable)
	}

	<-cancelCtx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	go func() {
		wg.Wait()
		shutdownCancel()
	}()
	<-shutdownCtx.Done()

	errorMessage := ""
	if shutdownCtx.Err() == context.DeadlineExceeded {
		for k, v := range states {
			if !v.done {
				errorMessage += fmt.Sprintf("timeout[%s] ", k)
			}
		}
	}

	for k, v := range states {
		if v.err != nil {
			errorMessage += fmt.Sprintf("error[%s: %v] ", k, v.err)
		}
	}

	if len(errorMessage) > 0 {
		return errors.New(errorMessage)
	}
	return nil
}
