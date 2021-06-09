package goroutine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Settings_setDefaults(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  Settings
		expected Settings
	}{
		"default settings": {
			expected: Settings{
				Timeout: time.Second,
			},
		},
		"all-set settings": {
			initial: Settings{
				Timeout:  time.Minute,
				Critical: true,
			},
			expected: Settings{
				Timeout:  time.Minute,
				Critical: true,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.initial.setDefaults()

			assert.Equal(t, testCase.expected, testCase.initial)
		})
	}
}
