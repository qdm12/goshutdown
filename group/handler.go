package group

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/handler"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Handler

// Handler handles a group of shutdown handlers.
type Handler interface {
	// Name returns the name set for this group handler.
	Name() string
	// IsCritical returs true if the group of goroutines handler is critical
	// and must be terminated before continuing other external shutdown procedures.
	IsCritical() bool
	// Shutdown initiates the shutdown process for all the goroutines of the group in parallel.
	// It executes onSuccess or onFailure if a goroutine completion is a success or a failure, respectively.
	// It returns the number of incomplete goroutines when done.
	Shutdown(ctx context.Context) (err error)
	// Add adds a goroutine to the group of goroutine handlers.
	Add(handlers ...handler.Handler)
}

type groupHandler struct {
	name     string
	settings Settings
	handlers []handler.Handler
}

func New(name string, settings Settings) Handler {
	settings.setDefaults()
	return &groupHandler{
		name:     name,
		settings: settings,
	}
}

func (h *groupHandler) Name() string {
	return h.name
}

func (h *groupHandler) IsCritical() bool {
	return h.settings.Critical
}

func (h *groupHandler) Add(handlers ...handler.Handler) {
	h.handlers = append(h.handlers, handlers...)
}

var (
	// ErrCriticalTimeout is the error when a critical goroutine shutdown timed out in the group.
	ErrCriticalTimeout = errors.New("critical shutdown timed out in the group")
	// ErrTimeout is the error when one of the group shutdown times out.
	ErrTimeout = errors.New("group shutdown timed out")
)

func (h *groupHandler) Shutdown(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type completionStatus struct {
		name     string
		critical bool
		err      error
	}
	completed := make(chan completionStatus)

	for _, handler := range h.handlers {
		go func(handler goroutine.Handler) {
			completed <- completionStatus{
				name:     handler.Name(),
				critical: handler.IsCritical(),
				err:      handler.Shutdown(ctx),
			}
		}(handler)
	}

	var criticalErr error
	var errorMessages []string
	for range h.handlers {
		status := <-completed
		if status.err == nil {
			h.settings.OnSuccess(status.name)
			continue
		}

		h.settings.OnFailure(status.name, status.err)

		if criticalErr == nil && status.critical {
			criticalErr = status.err
			cancel() // stop shutdown of other goroutines
		}

		if criticalErr == nil {
			errorMessages = append(errorMessages, status.name+": "+status.err.Error())
		}
	}

	if criticalErr != nil {
		return fmt.Errorf("%w: %s", ErrCriticalTimeout, criticalErr)
	}

	if len(errorMessages) == 0 {
		return nil
	}

	return fmt.Errorf("%w: %d out of %d goroutines: %s",
		ErrTimeout, len(errorMessages), len(h.handlers),
		strings.Join(errorMessages, ", "))
}
