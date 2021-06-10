package goshutdown

import (
	"context"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/group"
	"github.com/qdm12/goshutdown/order"
)

type (
	// GoRoutineSettings is a type alias for goroutine.Settings and
	// defines the shutdown settings for a goroutine handler.
	GoRoutineSettings = goroutine.Settings
	// GroupSettings is a type alias for group.Settings and
	// defines the shutdown settings for a group handler.
	GroupSettings = group.Settings
	// OrderSettings is a type alias for order.Settings and
	// defines the shutdown settings for an order handler.
	OrderSettings = order.Settings
)

// NewGoRoutineHandler creates a new handler for a goroutine using the
// name and settings given. It returns the handler as well as the context
// and signal done channel to use in the actual Goroutine.
func NewGoRoutineHandler(name string, settings goroutine.Settings) (
	h goroutine.Handler, ctx context.Context, done chan<- struct{}) {
	return goroutine.New(name, settings)
}

// NewGroupHandler creates a new group handler using the name and settings given.
func NewGroupHandler(name string, settings group.Settings) group.Handler {
	return group.New(name, settings)
}

// NewOrder creates a new order handler using the name and settings given.
func NewOrder(name string, settings order.Settings) order.Handler {
	return order.New(name, settings)
}
