package patterns

// Fanout is a function that takes a source channel of type T and a number n, and returns a slice of n channels that receive the same values as the source channel.
//
// Parameters:
// - source: a channel of type T that serves as the source of values.
// - n: an integer that represents the number of channels to create.
//
// Return:
// - []<-chan T: a slice of n channels that receive the same values as the source channel.
func Fanout[T any](source <-chan T, n int) []<-chan T {
	out := make([]<-chan T, 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan T)
		out = append(out, ch)

		go func() {
			defer close(ch)

			for val := range source {
				ch <- val // without just passing the value, we can do heavy work
			}
		}()
	}

	return out
}
