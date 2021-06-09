package order

import "time"

// Settings define configuration settings for the shutdown Order.
type Settings struct {
	// Timeout is the global timeout for all shutdown operations.
	// It defaults to 1s if left unset.
	Timeout time.Duration
	// Critical can be set to true to indicate the shutdown process should exit if
	// this order of shutdown handlers cannot be completed.
	Critical bool
	// OnSuccess defines a function to execute when an handler in the order
	// terminates successfully. It is disabled if it is left unset.
	OnSuccess func(name string)
	// OnSuccess defines a function to execute when an handler in the order
	// does not terminate on time. It is disabled if it is left unset.
	OnFailure func(name string, err error)
}

func (s *Settings) setDefaults() {
	if s.Timeout == 0 {
		s.Timeout = time.Second
	}

	if s.OnSuccess == nil {
		s.OnSuccess = defaultOnSuccess
	}

	if s.OnFailure == nil {
		s.OnFailure = defaultOnFailure
	}
}

func defaultOnSuccess(groupName string)            {}
func defaultOnFailure(groupName string, err error) {}
