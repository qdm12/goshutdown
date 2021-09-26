package group

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newSettings(t *testing.T) {
	t.Parallel()

	s := newSettings()

	expected := Settings{
		Timeout:   time.Second,
		OnSuccess: defaultOnSuccess,
		OnFailure: defaultOnFailure,
	}

	var errDummy = errors.New("dummy")

	assert.NotPanics(t, func() {
		s.OnSuccess("group")
		s.OnFailure("group", errDummy)
	})

	assertSettingsEqual(t, &expected, &s)
}

// asserts the Settings a and b are equal and clear the problematic fields
// that cannot be asserted without reflect such as functions.
func assertSettingsEqual(t *testing.T, a, b *Settings) {
	t.Helper()
	assert.Equal(t, reflect.ValueOf(a.OnFailure), reflect.ValueOf(b.OnFailure))
	a.OnFailure, b.OnFailure = nil, nil

	assert.Equal(t, reflect.ValueOf(a.OnSuccess), reflect.ValueOf(b.OnSuccess))
	a.OnSuccess, b.OnSuccess = nil, nil

	assert.Equal(t, a, b)
}
