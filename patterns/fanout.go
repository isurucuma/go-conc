package patterns

func Fanout[T any](source <-chan T, n int) []<-chan T {
	out := make([]<-chan T, 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan T)
		out = append(out, ch)

		go func() {
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}

	return out
}
