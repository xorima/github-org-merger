package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler_Gather(t *testing.T) {
	t.Run("It should return a list of repositories", func(t *testing.T) {
		h := NewHandler("")
		h.Gather()
		assert.False(t, true)
	})
}
