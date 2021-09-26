package group

import "time"

type Option func(s *Settings)

// OptionTimeout sets a timeout for the goroutine shutdown operation.
// Note the timeout defaults to one second.
func OptionTimeout(timeout time.Duration) Option {
	return func(s *Settings) {
		s.Timeout = timeout
	}
}

// OptionCritical marks the shutdown operation as critical.
func OptionCritical() Option {
	return func(s *Settings) {
		s.Critical = true
	}
}

// OptionOnSuccess sets a function to execute when the shutdown is a success.
func OptionOnSuccess(fn func(orderName string)) Option {
	return func(s *Settings) {
		s.OnSuccess = fn
	}
}

// OptionOnFailure sets a function to execute when the shutdown is a failure.
func OptionOnFailure(fn func(orderName string, err error)) Option {
	return func(s *Settings) {
		s.OnFailure = fn
	}
}
