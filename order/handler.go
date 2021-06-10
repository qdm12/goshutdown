package order

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/goshutdown/handler"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Handler

// Handler handles an order of shutdown handlers.
type Handler interface {
	// Name returns the name set for this order handler.
	Name() string
	// IsCritical returs true if the order handler is critical and must be terminated
	// before continuing other external shutdown procedures.
	IsCritical() bool
	// Shutdown initiates the shutdown process, starting with the start of the order
	// (group or single goroutine). It returns an error if one or more goroutines
	// did not complete on time, and nil otherwise. You can stop the shutdown process
	// by canceling its context, but really you should not do that.
	Shutdown(ctx context.Context) (err error)
	// Append appends one or more handlers to the order. An handler.Handler can be a
	// group.Handler, a goroutine.Handler or a user defined implementation.
	// The handlers are appended in a first-in-first-out fashion.
	Append(handlers ...handler.Handler)
}

type orderHandler struct {
	name     string
	settings Settings
	handlers []handler.Handler
}

// New creates a new shutdown Handler with the given settings.
// Each field of settings is set to its default if left unset.
func New(name string, settings Settings) Handler {
	settings.setDefaults()
	return &orderHandler{
		name:     name,
		settings: settings,
	}
}

func (h *orderHandler) Name() string {
	return h.name
}

func (h *orderHandler) IsCritical() bool {
	return h.settings.Critical
}

var (
	// ErrCriticalTimeout is the error when a critical shutdown timed out in the order.
	ErrCriticalTimeout = errors.New("critical order handler timed out")
	// ErrTimeout is the error when one or more shutdown timed out in the order.
	ErrTimeout = errors.New("ordered shutdown timed out")
)

func (h *orderHandler) Shutdown(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, h.settings.Timeout)
	defer cancel()

	var errorMessages []string //nolint:prealloc

	for _, handler := range h.handlers {
		name := handler.Name()

		err := handler.Shutdown(ctx)
		if err == nil {
			h.settings.OnSuccess(name)
			continue
		}

		h.settings.OnFailure(name, err)
		if handler.IsCritical() {
			return fmt.Errorf("%w: %s: %s", ErrCriticalTimeout, name, err)
		}
		errorMessages = append(errorMessages, name+": "+err.Error())
	}

	if len(errorMessages) == 0 {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrTimeout, strings.Join(errorMessages, "; "))
}

func (h *orderHandler) Append(handlers ...handler.Handler) {
	h.handlers = append(h.handlers, handlers...)
}
