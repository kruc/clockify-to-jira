package clockify

import (
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestInitClient(t *testing.T) {

	t.Run("Returns error on init client with invalid token", func(t *testing.T) {
		_, err := NewClient("")

		assert.Errors(t, err, ErrClockifyClientInitError)
	})
}

func TestError(t *testing.T) {
	t.Run("ErrNotFound", func(t *testing.T) {
		got := ClockifyErr("Error message").Error()
		want := "Error message"

		assert.Strings(t, got, want)
	})
}
