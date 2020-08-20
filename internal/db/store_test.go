package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseUUID(t *testing.T) {
	t.Run("Valid UUID", func(t *testing.T) {
		uid := "94cc5321-ec44-464f-9008-3d81f5e2c18f"

		result, err := ParseUUID(uid)
		assert.NoError(t, err, "Expected no error from parse")
		assert.Equal(t, uid, result.String(), "Expected string to be parsed to UUID")
	})

	t.Run("Empty UUID", func(t *testing.T) {
		result, err := ParseUUID("")
		assert.NoError(t, err, "Expected no error from parse")
		assert.Equal(t, uuid.Nil, result, "Expected string to be parsed to a 0 value UUID")
	})
}
