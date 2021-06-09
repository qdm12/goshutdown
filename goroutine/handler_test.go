package goroutine

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()

	const name = "routine name"
	const timeout = time.Second
	settings := Settings{
		Timeout:  timeout,
		Critical: true,
	}

	intf, ctx, done := New(name, settings)

	assert.NotNil(t, ctx)
	assert.NotNil(t, done)

	impl, ok := intf.(*handler)
	require.True(t, ok)

	assert.Equal(t, name, impl.name)
	assert.Equal(t, settings, impl.settings)
	// cannot assert cancel and done as they are hidden away.
}

func Test_handler_Name(t *testing.T) {
	t.Parallel()

	const name = "routine name"
	h := &handler{
		name: name,
	}
	s := h.Name()
	assert.Equal(t, name, s)
}

func Test_handler_IsCritical(t *testing.T) {
	t.Parallel()
	const critical = true

	h := &handler{
		settings: Settings{Critical: critical},
	}
	c := h.IsCritical()

	assert.Equal(t, critical, c)
}

func Test_handler_Shutdown(t *testing.T) {
	t.Parallel()

	t.Run("goroutine completes", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		done := make(chan struct{})
		close(done)

		h := &handler{
			cancel: func() {},
			done:   done,
			settings: Settings{
				Timeout: time.Hour,
			},
		}

		err := h.Shutdown(ctx)

		assert.NoError(t, err)
	})

	t.Run("shutdown context canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		h := &handler{
			cancel: func() {},
			done:   nil,
			settings: Settings{
				Timeout: time.Hour,
			},
		}

		err := h.Shutdown(ctx)

		require.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		h := &handler{
			cancel: func() {},
			done:   nil,
			settings: Settings{
				Timeout: time.Nanosecond,
			},
		}

		err := h.Shutdown(ctx)

		require.Error(t, err)
		expectedErr := errors.New("goroutine shutdown timed out: after 1ns")
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("completes with no timeout", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		done := make(chan struct{})
		close(done)

		h := &handler{
			cancel: func() {},
			done:   done,
			settings: Settings{
				Timeout: time.Hour,
			},
		}

		err := h.Shutdown(ctx)

		assert.NoError(t, err)
	})
}
