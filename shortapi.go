package goshutdown

import (
	"context"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/group"
	"github.com/qdm12/goshutdown/order"
)

type (
	GoRoutineSettings = goroutine.Settings
	GroupSettings     = group.Settings
	OrderSettings     = order.Settings
)

func NewGoRoutineHandler(name string, settings goroutine.Settings) (
	h goroutine.Handler, ctx context.Context, done chan<- struct{}) {
	return goroutine.New(name, settings)
}

func NewGroupHandler(name string, settings group.Settings) group.Handler {
	return group.New(name, settings)
}

func NewOrder(name string, settings order.Settings) order.Handler {
	return order.New(name, settings)
}
