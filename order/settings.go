package order

import "time"

// settings defines configuration settings for the shutdown Order.
type settings struct {
	// timeout is the global timeout for all shutdown operations.
	// It defaults to 1s if left unset.
	timeout time.Duration
	// critical can be set to true to indicate the shutdown process should exit if
	// this order of shutdown handlers cannot be completed.
	critical bool
	// onSuccess defines a function to execute when an handler in the order
	// terminates successfully. It is disabled if it is left unset.
	onSuccess func(name string)
	// OnSuccess defines a function to execute when an handler in the order
	// does not terminate on time. It is disabled if it is left unset.
	onFailure func(name string, err error)
}

func newSettings() settings {
	return settings{
		timeout:   time.Second,
		onSuccess: defaultOnSuccess,
		onFailure: defaultOnFailure,
	}
}

func defaultOnSuccess(name string)            {}
func defaultOnFailure(name string, err error) {}
