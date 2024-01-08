package cache

import (
	"context"
	"github.com/lisgo88/faraway-test/internal/config"
	"github.com/lisgo88/faraway-test/internal/pkg/repository"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

const clearInterval = 60 * time.Second

// Cache - generic key-value cache with time expiration
type Cache struct {
	cache map[string]value

	mu     sync.RWMutex
	config config.CacheClient
	logger zerolog.Logger
}

func New(ctx context.Context, cfg config.CacheClient, log zerolog.Logger) repository.Cache {
	cache := &Cache{
		cache:  make(map[string]value),
		config: cfg,
		logger: log,
	}

	// start cleaner worker
	go cache.clearExpiredWorker(ctx)

	return cache
}

// Get - get value by key
func (c *Cache) Get(key string) (val string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.cache[key]

	if ok && value.valid() {
		val, ok = value.data, true
		return
	}

	return val, false
}

// Set - add new value by key
func (c *Cache) Set(key, val string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = value{
		data: val,
		exp:  time.Now().Add(c.config.TTL * time.Second).UnixNano(), // set default ttl
	}

	return nil
}

// SetWithTTL - add new value by key with TTL
func (c *Cache) SetWithTTL(key, val string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = value{
		data: val,
		exp:  time.Now().Add(ttl * time.Second).UnixNano(),
	}

	return nil
}

// Delete - delete value by key
func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)

	return nil
}

// clearExpiredWorker - clear expired keys
func (c *Cache) clearExpiredWorker(ctx context.Context) {
	ticker := time.NewTicker(clearInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug().Any("cache", "clearExpiredWorker").Msg("context canceled")
			return
		case <-ticker.C:
			for key, val := range c.cache {
				if !val.valid() {
					_ = c.Delete(key)
				}
			}

			c.logger.Debug().Any("cache", "clearExpiredWorker").Msg("cache clear")
		}
	}
}

type value struct {
	data string
	exp  int64
}

func (v value) valid() bool {
	return v.exp == 0 || time.Now().UnixNano() < v.exp
}
