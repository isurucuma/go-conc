package patterns

import (
	"sync"
)

func Fanin[T any](channels []<-chan T) <-chan T {
	out := make(chan T)

	wg := &sync.WaitGroup{}
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan T) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
