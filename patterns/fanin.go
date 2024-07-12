package patterns

import (
	"sync"
)

// Fanin is a function that merges multiple input channels into a single output channel.
//
// It takes a slice of input channels 'channels' and returns a single output channel.
func Fanin[T any](channels []<-chan T) <-chan T {
	out := make(chan T)

	wg := &sync.WaitGroup{}
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan T) {
			defer wg.Done()
			for val := range c {
				out <- val // here instead of just passing, we can do a heavy work
			}
		}(ch)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
