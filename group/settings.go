package group

import "time"

// Settings define configuration settings for the shutdown Group.
type Settings struct {
	// Timeout is the timeout for termninating all the goroutines in the group.
	// It defaults to 1s if left unset.
	Timeout time.Duration
	// Critical can be set to true to indicate the shutdown process should exit if
	// this group of goroutines cannot be completed.
	Critical bool
	// OnSuccess defines a function to execute when a one of the goroutines
	// terminates successfully. It is disabled if it is left unset.
	OnSuccess func(goRoutineName string)
	// OnFailure defines a function to execute when a one of the goroutines
	// does not terminate on time. It is disabled if it is left unset.
	OnFailure func(goRoutineName string, err error)
}

func newSettings() Settings {
	return Settings{
		Timeout:   time.Second,
		OnSuccess: defaultOnSuccess,
		OnFailure: defaultOnFailure,
	}
}

func defaultOnSuccess(goRoutineName string)            {}
func defaultOnFailure(goRoutineName string, err error) {}
