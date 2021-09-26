package order

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newSettings(t *testing.T) {
	t.Parallel()

	var (
		errDummy = errors.New("dummy")
	)

	s := newSettings()

	assert.NotPanics(t, func() {
		s.onSuccess("group")
		s.onFailure("group", errDummy)
	})

	expected := settings{
		timeout:   time.Second,
		onSuccess: defaultOnSuccess,
		onFailure: defaultOnFailure,
	}

	assertSettingsEqual(t, &expected, &s)
}

// asserts the Settings a and b are equal and clear the problematic fields
// that cannot be asserted without reflect such as functions.
func assertSettingsEqual(t *testing.T, a, b *settings) {
	t.Helper()
	assert.Equal(t, reflect.ValueOf(a.onFailure), reflect.ValueOf(b.onFailure))
	a.onFailure, b.onFailure = nil, nil

	assert.Equal(t, reflect.ValueOf(a.onSuccess), reflect.ValueOf(b.onSuccess))
	a.onSuccess, b.onSuccess = nil, nil

	assert.Equal(t, a, b)
}
