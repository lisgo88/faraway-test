package hashcash

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseHashcash(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "clientID", "1234")

	t.Run("zero bits error", func(t *testing.T) {
		_, err := New(ctx, 0)
		require.EqualError(t, err, "zero bits must be more than zero")
	})

	t.Run("create success", func(t *testing.T) {
		original, err := New(ctx, 20)
		require.NoError(t, err)

		parsed, err := ParseHeader(string(original.Header()))
		require.NoError(t, err)
		require.Equal(t, original, parsed)

		parsed.counter++
		require.Equal(t, original.bits, parsed.bits)
	})

	t.Run("max attempts error", func(t *testing.T) {
		header := "1:7:20240108180137:resource::Cxphfw==:MA=="

		hashcash, err := ParseHeader(header)
		require.NoError(t, err)

		err = hashcash.Compute(100)
		require.EqualError(t, err, "max attempts exceeded")
	})

	t.Run("compute success", func(t *testing.T) {
		header := "1:3:20240108180137:resource::Cxphfw==:MA=="

		hashcash, err := ParseHeader(header)
		require.NoError(t, err)

		err = hashcash.Compute(10000)
		require.NoError(t, err)
		require.Equal(t, 5581, hashcash.Counter())
	})
}
