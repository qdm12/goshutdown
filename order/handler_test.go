package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/goroutine/mock_goroutine"
	"github.com/qdm12/goshutdown/handler"
	"github.com/qdm12/goshutdown/handler/mock_handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()

	const name = "name"

	expected := &orderHandler{
		name: name,
		settings: settings{
			timeout:   time.Second,
			onSuccess: defaultOnSuccess,
			onFailure: defaultOnFailure,
		},
	}

	intf := New(name)

	impl, ok := intf.(*orderHandler)
	require.True(t, ok)

	assertSettingsEqual(t, &expected.settings, &impl.settings)

	assert.Equal(t, expected, impl)
}

func Test_orderHandler_Name(t *testing.T) {
	t.Parallel()
	const name = "name"

	h := &orderHandler{
		name: name,
	}
	n := h.Name()

	assert.Equal(t, name, n)
}

func Test_orderHandler_IsCritical(t *testing.T) {
	t.Parallel()
	const critical = true

	h := &orderHandler{
		settings: settings{critical: critical},
	}
	c := h.IsCritical()

	assert.Equal(t, critical, c)
}

func Test_orderHandler_Shutdown(t *testing.T) {
	t.Parallel()

	type handlerReturnValues struct {
		name     string
		critical bool
		err      error
	}

	testCases := map[string]struct {
		o                    Handler
		handlersReturnValues []handlerReturnValues // dynamically set handlers in o in subtests
		err                  error
	}{
		"no handler": {
			o: New("order name"),
		},
		"single handler complete": {
			o:                    New("order name"),
			handlersReturnValues: []handlerReturnValues{{}},
		},
		"single handler failed": {
			o: New("order name"),
			handlersReturnValues: []handlerReturnValues{
				{name: "name", err: goroutine.ErrTimeout},
			},
			err: errors.New("ordered shutdown timed out: name: goroutine shutdown timed out"),
		},
		"two handlers complete": {
			o:                    New("order name"),
			handlersReturnValues: []handlerReturnValues{{}, {}},
		},
		"two handlers with one failed": {
			o: New("order name"),
			handlersReturnValues: []handlerReturnValues{
				{name: "name", err: goroutine.ErrTimeout},
				{},
			},
			err: errors.New("ordered shutdown timed out: name: goroutine shutdown timed out"),
		},
		"two handlers with first critical failed": {
			o: New("order name"),
			handlersReturnValues: []handlerReturnValues{
				{name: "name", critical: true, err: goroutine.ErrTimeout},
				{},
			},
			err: errors.New("critical order handler timed out: name: goroutine shutdown timed out"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			o := testCase.o

			ctx := context.Background()

			criticalFound := false
			for _, returnValues := range testCase.handlersReturnValues {
				handler := mock_handler.NewMockHandler(ctrl)
				if !criticalFound {
					handler.EXPECT().Name().Return(returnValues.name)
					handler.EXPECT().Shutdown(gomock.Any()).Return(returnValues.err)
					if returnValues.err != nil {
						handler.EXPECT().IsCritical().Return(returnValues.critical)
					}
					criticalFound = criticalFound || returnValues.critical
				}
				o.Append(handler)
			}

			err := testCase.o.Shutdown(ctx)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_orderHandler_Append(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	a := mock_goroutine.NewMockHandler(ctrl)
	b := mock_goroutine.NewMockHandler(ctrl)

	o := &orderHandler{
		handlers: []handler.Handler{a},
	}

	o.Append(b)

	expectedHandler := &orderHandler{
		handlers: []handler.Handler{a, b},
	}

	assert.Equal(t, expectedHandler, o)
}
