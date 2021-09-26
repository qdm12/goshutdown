package goroutine

import "time"

type Option func(s *settings)

// OptionTimeout sets a timeout for the goroutine shutdown operation.
// Note the timeout defaults to one second.
func OptionTimeout(timeout time.Duration) Option {
	return func(s *settings) {
		s.timeout = timeout
	}
}

// OptionCritical marks the shutdown operation as critical.
func OptionCritical() Option {
	return func(s *settings) {
		s.critical = true
	}
}
