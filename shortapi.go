package goshutdown

import (
	"context"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/group"
	"github.com/qdm12/goshutdown/order"
)

type (
	// GoRoutineOption is a type alias for goroutine.Option and
	// defines the shutdown options for a goroutine handler.
	GoRoutineOption = goroutine.Option
	// GroupOption is a type alias for group.Option and
	// defines the shutdown options for a group handler.
	GroupOption = group.Option
	// OrderOptions is a type alias for order.Option and
	// defines the shutdown options for an order handler.
	OrderOptions = order.Option
)

// NewGoRoutineHandler creates a new handler for a goroutine using the
// name and options given. It returns the handler as well as the context
// and signal done channel to use in the actual Goroutine.
func NewGoRoutineHandler(name string, options ...goroutine.Option) (
	h goroutine.Handler, ctx context.Context, done chan<- struct{}) {
	return goroutine.New(name, options...)
}

// NewGroupHandler creates a new group handler using the name and options given.
func NewGroupHandler(name string, options ...group.Option) group.Handler {
	return group.New(name, options...)
}

// NewOrderHandler creates a new order handler using the name and options given.
func NewOrderHandler(name string, options ...order.Option) order.Handler {
	return order.New(name, options...)
}
