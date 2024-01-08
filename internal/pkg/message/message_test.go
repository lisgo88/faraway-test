package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lisgo88/faraway-test/internal/pkg/message"
)

func TestParseMessage(t *testing.T) {
	t.Run("incorrect format error", func(t *testing.T) {
		_, err := message.ParseMessage("1\n")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "incorrect message format")
	})

	t.Run("error success", func(t *testing.T) {
		msg, err := message.ParseMessage("0:payload\n")

		assert.NoError(t, err)
		assert.Equal(t, message.CommandError, msg.Command)
		assert.Equal(t, "payload", msg.Payload)
	})

	t.Run("request challenge success", func(t *testing.T) {
		msg, err := message.ParseMessage("1:payload\n")

		assert.NoError(t, err)
		assert.Equal(t, message.CommandRequestChallenge, msg.Command)
		assert.Equal(t, "payload", msg.Payload)
	})

	t.Run("response challenge success", func(t *testing.T) {
		msg, err := message.ParseMessage("2:payload\n")

		assert.NoError(t, err)
		assert.Equal(t, message.CommandResponseChallenge, msg.Command)
		assert.Equal(t, "payload", msg.Payload)
	})

	t.Run("request resource success", func(t *testing.T) {
		msg, err := message.ParseMessage("3:payload\n")

		assert.NoError(t, err)
		assert.Equal(t, message.CommandRequestResource, msg.Command)
		assert.Equal(t, "payload", msg.Payload)
	})

	t.Run("response resource success", func(t *testing.T) {
		msg, err := message.ParseMessage("4:payload\n")

		assert.NoError(t, err)
		assert.Equal(t, message.CommandResponseResource, msg.Command)
		assert.Equal(t, "payload", msg.Payload)
	})

	t.Run("unknown command error", func(t *testing.T) {
		_, err := message.ParseMessage("5:payload\n")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "unknown command")
	})
}
