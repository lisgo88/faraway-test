package loadenv

import (
	"os"
	"strconv"
	"time"
)

func String(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func Int(key string, defaultValue int) int {
	if value, err := strconv.Atoi(String(key, "")); err == nil {
		return value
	}

	return defaultValue
}

func Duration(key string, defaultValue time.Duration) time.Duration {
	if value, err := strconv.Atoi(String(key, "")); err == nil {
		return time.Duration(value)
	}

	return defaultValue
}
