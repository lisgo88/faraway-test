package app

import (
	"bufio"
	"context"
	"fmt"
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
	defer conn.Close()

	// client will send new request every 1 second endlessly
	for {
		msg, err := c.handleConnection(ctx, conn, conn)
		if err != nil {
			return err
		}

		c.logger.Info().Msgf("got result quote: %s", msg)
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) handleConnection(_ context.Context, readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	// 1 - request challenge
	if err := c.writeMsg(message.Message{Command: message.CommandRequestChallenge}, writerConn); err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	// read response
	rawMsg, err := c.readMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}

	// parse raw message
	msg, err := message.ParseMessage(rawMsg)
	if err != nil {
		c.logger.Error().Err(err).Any("raw msg", rawMsg).Msg("err parse msg")
		return "", fmt.Errorf("err parse msg: %w", err)
	}

	hc, err := hashcash.ParseHeader(msg.Payload)
	if err != nil {
		return "", fmt.Errorf("err parse hashcash: %w", err)
	}

	// 2. got challenge, compute hashcash
	if err := hc.Compute(c.config.MaxAttempts); err != nil {
		c.logger.Debug().Any("hashcash", hc).Msg("err compute hashcash")
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}

	// 3 - request resource with challenge
	if err := c.writeMsg(message.Message{Command: message.CommandRequestResource, Payload: string(hc.Header())}, writerConn); err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	// 4 - response with resource
	rawMsg, err = c.readMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}

	// parse raw message
	msg, err = message.ParseMessage(rawMsg)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}

	return msg.Payload, nil
}

func (c *Client) readMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString(message.DelimiterMessage)
}

func (c *Client) writeMsg(msg message.Message, w io.Writer) error {
	if _, err := w.Write(msg.Bytes()); err != nil {
		c.logger.Error().Err(err).Msg("io write error")
		return err
	}

	return nil
}
