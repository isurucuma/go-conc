package patterns

import (
	"context"
	"sync"
	"time"
)

type Future[T any] interface {
	Result() (T, error)
}

type InnerFuture[T any] struct {
	once sync.Once
	wg   sync.WaitGroup

	res   T
	err   error
	resCh <-chan T
	errCh <-chan error
}

func (f *InnerFuture[T]) Result() (T, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()

		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	f.wg.Wait()

	return f.res, f.err
}

func SlowFunction(ctx context.Context) Future[string] {
	resCh := make(chan string)
	errCh := make(chan error)
	go func() {
		select {
		case <-time.After(time.Second * 2):
			resCh <- "SlowFunc Execution done"
			errCh <- nil
		case <-ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	return &InnerFuture[string]{
		resCh: resCh,
		errCh: errCh,
	}
}
