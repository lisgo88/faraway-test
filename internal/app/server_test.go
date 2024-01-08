package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/lisgo88/faraway-test/internal/config"
	"github.com/lisgo88/faraway-test/internal/pkg/message"
	"github.com/lisgo88/faraway-test/internal/pkg/repository/mock"
)

func TestNewServer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCache := mock.NewMockCache(ctrl)
		mockQuotes := mock.NewMockQuotes(ctrl)

		cfg := config.ServerConfig{
			Host:        "localhost",
			Port:        8080,
			MaxAttempts: 100,
			HashBits:    3,
			HashTTL:     1000,
			TimeOut:     10,
		}

		req := fmt.Sprintf("%v:%s", message.CommandRequestChallenge, "payload")
		server := NewServer(mockCache, mockQuotes, cfg, zerolog.Logger{})
		mockCache.EXPECT().SetWithTTL(gomock.Any(), gomock.Any(), cfg.HashTTL).Return(nil)

		msg, err := server.processRequest(context.Background(), req, "clientInfo")

		assert.NoError(t, err)
		assert.NotNil(t, msg)
	})
}
