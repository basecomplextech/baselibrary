package async

// Routines groups multiple routines into a slice.
type Routines[T any] []Routine[T]

// Stop requests all routines to stop.
func (rr Routines[T]) Stop() {
	for _, r := range rr {
		r.Stop()
	}
}

// Wait starts a routine which waits for all routines and returns their results.
func (rr Routines[T]) Wait() Routine[[]Result[T]] {
	local := make(Routines[T], len(rr))
	copy(local, rr)

	return Call(func(stop <-chan struct{}) ([]Result[T], error) {
		results := make([]Result[T], 0, len(local))

		for _, r := range local {
			select {
			case <-r.Wait():
			case <-stop:
				return results, Stopped
			}

			v, err := r.Result()
			result := Result[T]{
				Err:   err,
				Value: v,
			}
			results = append(results, result)
		}

		return results, nil
	})
}
