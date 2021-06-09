package group

import (
	"context"
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

	const name = "group name"
	settings := Settings{
		Timeout:   time.Second,
		OnSuccess: defaultOnSuccess,
		OnFailure: defaultOnFailure,
	}

	expected := &groupHandler{
		name:     name,
		settings: settings,
	}

	intf := New(name, settings)

	impl, ok := intf.(*groupHandler)
	require.True(t, ok)

	assertSettingsEqual(t, &expected.settings, &impl.settings)

	assert.Equal(t, expected, impl)
}

func Test_groupHandler_Name(t *testing.T) {
	t.Parallel()
	const name = "name"

	h := &groupHandler{
		name: name,
	}
	n := h.Name()

	assert.Equal(t, name, n)
}

func Test_groupHandler_IsCritical(t *testing.T) {
	t.Parallel()
	const critical = true

	h := &groupHandler{
		settings: Settings{Critical: critical},
	}
	c := h.IsCritical()

	assert.Equal(t, critical, c)
}

func Test_groupHandler_Add(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	h := new(groupHandler)
	mockHandler := mock_handler.NewMockHandler(ctrl)

	h.Add(mockHandler)

	expectedHandler := &groupHandler{
		handlers: []handler.Handler{mockHandler},
	}

	assert.Equal(t, expectedHandler, h)
}

func Test_groupHandler_Shutdown_success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	const testGroupName = "group"

	const goRoutine1Name = "my-completed"
	goRoutine1 := mock_goroutine.NewMockHandler(ctrl)
	goRoutine1.EXPECT().Name().Return(goRoutine1Name)
	goRoutine1.EXPECT().IsCritical().Return(true)
	goRoutine1.EXPECT().Shutdown(gomock.Any()).Return(nil)

	const goRoutine2Name = "my-timed-out"
	goRoutine2 := mock_goroutine.NewMockHandler(ctrl)
	goRoutine2.EXPECT().Name().Return(goRoutine2Name)
	goRoutine2.EXPECT().IsCritical().Return(true)
	goRoutine2.EXPECT().Shutdown(gomock.Any()).Return(nil)

	onSuccess := func(goroutineName string) {
		switch goroutineName {
		case goRoutine1Name, goRoutine2Name:
		default:
			t.Errorf("onSuccess goRoutineName %q is not expected", goroutineName)
		}
	}
	onFailure := func(goroutineName string, err error) {
		t.Errorf("onFailure should not be called for %q with error: %s",
			goroutineName, err)
	}
	settings := Settings{
		OnSuccess: onSuccess,
		OnFailure: onFailure,
	}

	h := &groupHandler{
		name:     testGroupName,
		settings: settings,
		handlers: []handler.Handler{
			goRoutine1,
			goRoutine2,
		},
	}

	ctx := context.Background()

	err := h.Shutdown(ctx)

	assert.NoError(t, err)
}

func Test_groupHandler_Shutdown_one_timeout(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	const testGroupName = "group"

	const goRoutineSuccessName = "my-completed"
	goRoutineSuccess := mock_goroutine.NewMockHandler(ctrl)
	goRoutineSuccess.EXPECT().Name().Return(goRoutineSuccessName)
	goRoutineSuccess.EXPECT().IsCritical().Return(false)
	goRoutineSuccess.EXPECT().Shutdown(gomock.Any()).Return(nil)

	const goRoutineTimedoutName = "my-timed-out"
	goRoutineTimedout := mock_goroutine.NewMockHandler(ctrl)
	goRoutineTimedout.EXPECT().Name().Return(goRoutineTimedoutName)
	goRoutineTimedout.EXPECT().IsCritical().Return(false)
	goRoutineTimedout.EXPECT().Shutdown(gomock.Any()).Return(goroutine.ErrTimeout)

	onSuccess := func(goroutineName string) {
		assert.Equal(t, goRoutineSuccessName, goroutineName)
	}
	onFailure := func(goroutineName string, err error) {
		assert.Equal(t, goRoutineTimedoutName, goroutineName)
		assert.Equal(t, goroutine.ErrTimeout, err)
	}
	settings := Settings{
		OnSuccess: onSuccess,
		OnFailure: onFailure,
	}

	h := &groupHandler{
		name:     testGroupName,
		settings: settings,
		handlers: []handler.Handler{
			goRoutineSuccess,
			goRoutineTimedout,
		},
	}

	ctx := context.Background()

	err := h.Shutdown(ctx)

	require.Error(t, err)
	const expectedErrMessage = "group shutdown timed out: 1 out of 2 goroutines: my-timed-out: goroutine shutdown timed out" //nolint:lll
	assert.Equal(t, expectedErrMessage, err.Error())
}

func Test_groupHandler_Shutdown_critical_timeout(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	const testGroupName = "group"

	const goRoutineCriticalName = "my-critical"
	goRoutineCritical := mock_goroutine.NewMockHandler(ctrl)
	goRoutineCritical.EXPECT().Name().Return(goRoutineCriticalName)
	const critical = true
	goRoutineCritical.EXPECT().IsCritical().Return(critical)
	goRoutineCritical.EXPECT().Shutdown(gomock.Any()).Return(goroutine.ErrTimeout)

	const goRoutineTimedoutName = "my-timed-out"
	goRoutineTimedout := mock_goroutine.NewMockHandler(ctrl)
	goRoutineTimedout.EXPECT().Name().Return(goRoutineTimedoutName)
	goRoutineTimedout.EXPECT().IsCritical().Return(false)
	goRoutineTimedout.EXPECT().Shutdown(gomock.Any()).Return(goroutine.ErrTimeout)

	onSuccess := func(goroutineName string) {
		t.Errorf("onSuccess should not be called for %q", goroutineName)
	}
	onFailure := func(goroutineName string, err error) {
		switch goroutineName {
		case goRoutineCriticalName, goRoutineTimedoutName:
		default:
			t.Errorf("onFailure goRoutineName %q is not expected", goroutineName)
		}
		assert.Equal(t, goroutine.ErrTimeout, err)
	}
	settings := Settings{
		OnSuccess: onSuccess,
		OnFailure: onFailure,
	}

	h := &groupHandler{
		name:     testGroupName,
		settings: settings,
		handlers: []handler.Handler{
			goRoutineCritical,
			goRoutineTimedout,
		},
	}

	ctx := context.Background()

	err := h.Shutdown(ctx)

	require.Error(t, err)
	const expectedErrMessage = "critical shutdown timed out in the group: goroutine shutdown timed out"
	assert.Equal(t, expectedErrMessage, err.Error())
}
