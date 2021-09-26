package goroutine

import "time"

// settings defines configuration settings for the shutdown GoRoutine.
type settings struct {
	// timeout is the timeout for terminating the goroutine.
	// It defaults to 1s if left unset.
	timeout time.Duration
	// critical can be set to true to indicate the shutdown process should exit if
	// this goroutine cannot be terminated.
	critical bool
}

func newSettings() settings {
	return settings{
		timeout: time.Second,
	}
}
