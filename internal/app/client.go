package app

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/lisgo88/faraway-test/internal/config"
	"github.com/lisgo88/faraway-test/internal/pkg/hashcash"
	"github.com/lisgo88/faraway-test/internal/pkg/message"
)

// Client - client struct
type Client struct {
	config config.ClientConfig
	logger zerolog.Logger
}

// NewClient - create new client
func NewClient(cfg config.ClientConfig, log zerolog.Logger) *Client {
	return &Client{
		config: cfg,
		logger: log,
	}
}

// Run - start client
func (c *Client) Run(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.logger.Info().Any("config", conn.LocalAddr()).Msg("TCP client starting...")
	defer func() {
		_ = conn.Close()
	}()

	// client will send new request every 1 second endlessly
	for {
		res, reqID, err := c.handleConnection(ctx, conn, conn)
		if err != nil {
			c.logger.Error().Err(err).Any("requestID", reqID).Msg("err handle connection")
			_ = conn.Close()

			return err
		}

		if res != "" {
			c.logger.Info().Any("requestID", reqID).Msgf("got result quote: %s", res)
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *Client) handleConnection(ctx context.Context, r io.Reader, w io.Writer) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	buffer := make([]byte, 1024)
	length, err := r.Read(buffer)
	if err != nil {
		c.logger.Error().Err(err).Msg("read from conn")
		return "", "", err
	}

	serverMsg := &message.Message{}
	if err := proto.Unmarshal(buffer[:length], serverMsg); err != nil {
		c.logger.Error().Err(err).Msg("unmarshal msg")
		return "", "", err
	}

	if serverMsg.Challenge != "" {
		hc, err := hashcash.ParseHeader(serverMsg.Challenge)
		if err != nil {
			c.logger.Error().Err(err).Msg("parse hashcash")
			return "", serverMsg.RequestID, err
		}

		if err := hc.Compute(c.config.MaxAttempts); err != nil {
			c.logger.Debug().Any("hashcash", hc).Msg("compute hashcash")
			return "", serverMsg.RequestID, err
		}

		msg := &message.Message{
			RequestID: serverMsg.RequestID,
			ClientID:  serverMsg.ClientID,
			Challenge: string(hc.Header()),
		}

		clientRes, err := proto.Marshal(msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("marshal msg")
			return "", serverMsg.RequestID, err
		}

		if err := c.writeMsg(clientRes, w); err != nil {
			return "", serverMsg.RequestID, err
		}
	}

	return serverMsg.Payload, serverMsg.RequestID, nil
}

func (c *Client) writeMsg(msg []byte, w io.Writer) error {
	if _, err := w.Write(msg); err != nil {
		c.logger.Error().Err(err).Msg("io write error")
		return err
	}

	return nil
}
