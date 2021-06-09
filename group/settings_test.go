package group

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Settings_setDefaults(t *testing.T) {
	t.Parallel()

	var (
		onSuccess = func(groupName string) { _ = groupName }
		onFailure = func(groupName string, err error) { _ = groupName }
	)

	testCases := map[string]struct {
		initial  Settings
		expected Settings
	}{
		"default settings": {
			expected: Settings{
				Timeout:   time.Second,
				OnSuccess: defaultOnSuccess,
				OnFailure: defaultOnFailure,
			},
		},
		"all-set settings": {
			initial: Settings{
				Timeout:   time.Minute,
				OnSuccess: onSuccess,
				OnFailure: onFailure,
			},
			expected: Settings{
				Timeout:   time.Minute,
				OnSuccess: onSuccess,
				OnFailure: onFailure,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.initial.setDefaults()

			var errDummy = errors.New("dummy")

			assert.NotPanics(t, func() {
				testCase.initial.OnSuccess("group")
				testCase.initial.OnFailure("group", errDummy)
			})

			assertSettingsEqual(t, &testCase.expected, &testCase.initial)
		})
	}
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
