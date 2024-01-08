package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/lisgo88/faraway-test/internal/pkg/loadenv"
)

type Config struct {
	Env      string        `json:"env"`
	LogLevel zerolog.Level `json:"log_level"`

	Cache CacheClient `json:"cache"`

	Client ClientConfig `json:"client"`
	Server ServerConfig `json:"server"`
}

type CacheClient struct {
	TTL time.Duration `json:"ttl"`
}

type ClientConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	MaxAttempts int `json:"max_attempts"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	MaxAttempts int           `json:"max_attempts"`
	HashBits    int           `json:"hash_bits"`
	HashTTL     time.Duration `json:"hash_ttl"`
	TimeOut     time.Duration `json:"timeout"`
}

func GetConfig() (*Config, error) {
	var (
		filePath string
		configs  *Config
	)

	if os.Getenv("config") != "" {
		filePath = os.Getenv("config")
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		filePath = currentDir + "/internal/config/config.json"
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&configs)
	if err != nil {
		return nil, err
	}

	configs.Env = loadenv.String("ENV", configs.Env)

	configs.Cache.TTL = loadenv.Duration("CACHE_TTL", configs.Cache.TTL)

	configs.Client.Host = loadenv.String("CLIENT_HOST", configs.Client.Host)
	configs.Client.Port = loadenv.Int("CLIENT_PORT", configs.Client.Port)
	configs.Client.MaxAttempts = loadenv.Int("CLIENT_MAX_ATTEMPTS", configs.Client.MaxAttempts)

	configs.Server.Host = loadenv.String("SERVER_HOST", configs.Server.Host)
	configs.Server.Port = loadenv.Int("SERVER_PORT", configs.Server.Port)
	configs.Server.MaxAttempts = loadenv.Int("SERVER_MAX_ATTEMPTS", configs.Server.MaxAttempts)
	configs.Server.HashBits = loadenv.Int("SERVER_HASH_BITS", configs.Server.HashBits)
	configs.Server.HashTTL = loadenv.Duration("SERVER_HASH_TTL", configs.Server.HashTTL)
	configs.Server.TimeOut = loadenv.Duration("SERVER_TIMEOUT", configs.Server.TimeOut)

	return configs, nil
}
