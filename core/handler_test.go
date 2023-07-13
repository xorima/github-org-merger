package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Run("It should return a handler which is not nil", func(t *testing.T) {
		got := NewHandler("token")
		assert.NotNil(t, got.client)
	})
}
