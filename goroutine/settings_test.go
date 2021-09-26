package goroutine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newSettings(t *testing.T) {
	t.Parallel()

	s := newSettings()

	expected := settings{
		timeout: time.Second,
	}

	assert.Equal(t, expected, s)
}
