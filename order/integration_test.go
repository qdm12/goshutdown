package order

import (
	"context"
	"testing"
	"time"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func functionA(ctx context.Context, done chan<- struct{}) {
	<-ctx.Done()
	close(done)
}

func functionB(ctx context.Context, _ chan<- struct{}) {
	<-ctx.Done()
}

func Test_Handler_NoHandlers(t *testing.T) {
	t.Parallel()
	order := New("order")
	err := order.Shutdown(context.Background())
	require.NoError(t, err)
}

func Test_Handler_GoRoutines_FirstFails(t *testing.T) {
	t.Parallel()
	order := New("order", OptionTimeout(2*time.Second))

	handlerB, ctxB, doneB := goroutine.New("B", goroutine.OptionTimeout(time.Nanosecond))
	go functionB(ctxB, doneB)
	order.Append(handlerB)

	handlerA, ctxA, doneA := goroutine.New("A")
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	err := order.Shutdown(context.Background())
	require.Error(t, err)
	assert.Equal(t, "ordered shutdown timed out: B: goroutine shutdown timed out: after 1ns", err.Error())
}

func Test_Handler_GoRoutines_FirstFailsCritical(t *testing.T) {
	t.Parallel()
	order := New("order", OptionTimeout(2*time.Second))

	handlerB, ctxB, doneB := goroutine.New("B", goroutine.OptionTimeout(time.Nanosecond), goroutine.OptionCritical())
	go functionB(ctxB, doneB)
	order.Append(handlerB)

	handlerA, ctxA, doneA := goroutine.New("A")
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	err := order.Shutdown(context.Background())
	require.Error(t, err)
	assert.Equal(t, "critical order handler timed out: B: goroutine shutdown timed out: after 1ns", err.Error())
}

func Test_Handler_GoRoutines_SecondFails(t *testing.T) {
	t.Parallel()
	order := New("order", OptionTimeout(2*time.Second))

	handlerA, ctxA, doneA := goroutine.New("A")
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	handlerB, ctxB, doneB := goroutine.New("B", goroutine.OptionTimeout(time.Nanosecond))
	go functionB(ctxB, doneB)
	order.Append(handlerB)

	err := order.Shutdown(context.Background())
	require.Error(t, err)
	assert.Equal(t, "ordered shutdown timed out: B: goroutine shutdown timed out: after 1ns", err.Error())
}

func Test_Handler_GoRoutines_SecondFailsCritical(t *testing.T) {
	t.Parallel()
	order := New("order", OptionTimeout(2*time.Second))

	handlerA, ctxA, doneA := goroutine.New("A")
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	handlerB, ctxB, doneB := goroutine.New("B", goroutine.OptionTimeout(time.Nanosecond), goroutine.OptionCritical())
	go functionB(ctxB, doneB)
	order.Append(handlerB)

	err := order.Shutdown(context.Background())
	require.Error(t, err)
	assert.Equal(t, "critical order handler timed out: B: goroutine shutdown timed out: after 1ns", err.Error())
}
