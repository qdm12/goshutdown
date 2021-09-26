package goroutine

import (
	"context"
	"errors"
	"fmt"
	"time"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Handler

// Handler handles a goroutine shutdown handler.
type Handler interface {
	// Name returns the name set for this goroutine.
	Name() string
	// IsCritical returns true if the goroutine is critical and must be terminated.
	IsCritical() bool
	// Shutdown shuts a goroutine down by canceling its associated context.
	// It then waits for the goroutine to close its done signal channel.
	// If the shutdown context is done, it returns the context error.
	// If the goroutine specific timeout is reached, it returns a timeout error.
	// indicating the goroutine did not terminate.
	Shutdown(ctx context.Context) (err error)
}

// New creates a goroutine handler with a timeout if timeout > 0.
func New(name string, options ...Option) (
	h Handler, ctx context.Context, done chan<- struct{}) {
	settings := newSettings()
	for _, option := range options {
		option(&settings)
	}

	ctx, cancel := context.WithCancel(context.Background())
	bidirectionalDone := make(chan struct{})

	h = &handler{
		name:     name,
		settings: settings,
		cancel:   cancel,
		done:     bidirectionalDone,
	}

	return h, ctx, bidirectionalDone
}

type handler struct {
	name     string
	settings settings
	cancel   context.CancelFunc
	done     <-chan struct{}
}

func (h *handler) Name() string {
	return h.name
}

func (h *handler) IsCritical() bool {
	return h.settings.critical
}

// ErrTimeout is the error when the goroutine shutdown times out.
var ErrTimeout = errors.New("goroutine shutdown timed out")

func (h *handler) Shutdown(ctx context.Context) (err error) {
	timer := time.NewTimer(h.settings.timeout)
	if h.settings.timeout == 0 {
		timer.Stop()
	}

	h.cancel()

	select {
	case <-h.done:
		if h.settings.timeout > 0 && !timer.Stop() {
			<-timer.C
		}
		return nil
	case <-ctx.Done():
		if h.settings.timeout > 0 && !timer.Stop() {
			<-timer.C
		}
		return ctx.Err() //nolint:wrapcheck
	case <-timer.C:
		return fmt.Errorf("%w: after %s", ErrTimeout, h.settings.timeout)
	}
}
