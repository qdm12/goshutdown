package goroutine

import (
	"context"
	"errors"
	"fmt"
	"time"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Handler

type Handler interface {
	// Name returns the name set for this goroutine.
	Name() string
	// IsCritical returs true if the goroutine is critical and must be terminated.
	IsCritical() bool
	// Shutdown shuts a goroutine down by canceling its associated context.
	// It then waits for the goroutine to close its done signal channel.
	// If the shutdown context is done, it returns the context error.
	// If the goroutine specific timeout is reached, it returns a timeout error.
	// indicating the goroutine did not terminate.
	Shutdown(ctx context.Context) (err error)
}

// New creates a goroutine handler with a timeout if timeout > 0.
func New(name string, settings Settings) (
	h Handler, ctx context.Context, done chan<- struct{}) {
	settings.setDefaults()

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
	settings Settings
	cancel   context.CancelFunc
	done     <-chan struct{}
}

func (h *handler) Name() string {
	return h.name
}

func (h *handler) IsCritical() bool {
	return h.settings.Critical
}

// ErrTimeout is the error when the goroutine shutdown times out.
var ErrTimeout = errors.New("goroutine shutdown timed out")

func (h *handler) Shutdown(ctx context.Context) (err error) {
	timer := time.NewTimer(h.settings.Timeout)
	if h.settings.Timeout == 0 {
		timer.Stop()
	}

	h.cancel()

	select {
	case <-h.done:
		if h.settings.Timeout > 0 && !timer.Stop() {
			<-timer.C
		}
		return nil
	case <-ctx.Done():
		if h.settings.Timeout > 0 && !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case <-timer.C:
		return fmt.Errorf("%w: after %s", ErrTimeout, h.settings.Timeout)
	}
}
