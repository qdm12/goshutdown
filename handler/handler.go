package handler

import "context"

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Handler

// Handler is the minimal common interface for shutdown items.
type Handler interface {
	// Name returns the name assigned to the handler.
	Name() string
	// IsCritical returns true if the shutdown process is critical and further
	// operations should be dropped it it cannot be done.
	IsCritical() bool
	// Shutdown initiates the shutdown process and returns an error if it fails.
	Shutdown(ctx context.Context) (err error)
}
