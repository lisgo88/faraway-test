package repository

import (
	"time"
)

//go:generate mockgen -destination=mock/repository.go -package=mock_repository faraway-test/internal/pkg/repository Quotes,Cache
type Quotes interface {
	GetQuote() (string, error)
}

type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string) error
	SetWithTTL(key, value string, ttl time.Duration) error
	Delete(key string) error
}
