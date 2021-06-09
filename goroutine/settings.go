package goroutine

import "time"

// Settings define configuration settings for the shutdown GoRoutine.
type Settings struct {
	// Timeout is the timeout for terminating the goroutine.
	// It defaults to 1s if left unset.
	Timeout time.Duration
	// Critical can be set to true to indicate the shutdown process should exit if
	// this goroutine cannot be terminated.
	Critical bool
}

func (s *Settings) setDefaults() {
	if s.Timeout == 0 {
		s.Timeout = time.Second
	}
}
