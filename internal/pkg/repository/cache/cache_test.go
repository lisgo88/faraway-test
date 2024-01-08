package cache

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/lisgo88/faraway-test/internal/config"
)

func TestCache(t *testing.T) {
	t.Run("get error", func(t *testing.T) {
		cfg := config.CacheClient{
			TTL: 10,
		}

		cacheRepo := New(context.Background(), cfg, zerolog.Logger{})
		_, ok := cacheRepo.Get("key")
		assert.Equal(t, ok, false)
	})

	t.Run("get success", func(t *testing.T) {
		cfg := config.CacheClient{
			TTL: 10,
		}

		cacheRepo := New(context.Background(), cfg, zerolog.Logger{})

		err := cacheRepo.Set("key", "value")
		assert.NoError(t, err)

		val, ok := cacheRepo.Get("key")
		assert.Equal(t, ok, true)
		assert.Equal(t, val, "value")
	})

	t.Run("delete success", func(t *testing.T) {
		cfg := config.CacheClient{
			TTL: 10,
		}

		cacheRepo := New(context.Background(), cfg, zerolog.Logger{})

		err := cacheRepo.Set("key", "value")
		assert.NoError(t, err)

		val, ok := cacheRepo.Get("key")
		assert.Equal(t, ok, true)
		assert.Equal(t, val, "value")

		err = cacheRepo.Delete("key")
		assert.NoError(t, err)

		_, ok = cacheRepo.Get("key")
		assert.Equal(t, ok, false)
	})
}
