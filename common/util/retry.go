package util

// Retry calls handler multiple times.
func Retry(handler func() error, times int) error {
	var err error
	for i := 0; i < times; i++ {
		err = handler()
		if err == nil {
			return nil
		}
	}
	return err
}
