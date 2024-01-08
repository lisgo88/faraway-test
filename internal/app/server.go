package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/lisgo88/faraway-test/internal/config"
	"github.com/lisgo88/faraway-test/internal/pkg/hashcash"
	"github.com/lisgo88/faraway-test/internal/pkg/message"
	"github.com/lisgo88/faraway-test/internal/pkg/repository"
)

type Server struct {
	cacheRepo  repository.Cache
	quotesRepo repository.Quotes

	config config.ServerConfig
	logger zerolog.Logger
}

func NewServer(cr repository.Cache, qr repository.Quotes, cfg config.ServerConfig, log zerolog.Logger) *Server {
	return &Server{
		cacheRepo:  cr,
		quotesRepo: qr,
		config:     cfg,
		logger:     log,
	}
}

// Run - start server
func (s *Server) Run(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer listener.Close()
	s.logger.Info().Any("config", listener.Addr()).Msg("TCP server starting...")

	for {
		// listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}

		// set timeout for connection
		_ = conn.SetReadDeadline(time.Now().Add(s.config.TimeOut * time.Second))

		// handle connections
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	s.logger.Info().Str("client", conn.RemoteAddr().String()).Msg("new client")
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		rawMsg, err := reader.ReadString(message.DelimiterMessage)
		if err != nil {
			s.logger.Error().Err(err).Any("rawMsg", rawMsg).Msg("read connection")
			s.writeMsg(message.Message{}, err, conn)

			return
		}

		msg, err := s.processRequest(ctx, rawMsg, conn.RemoteAddr().String())
		if err != nil {
			s.logger.Error().Err(err).Msg("request processing")
			s.writeMsg(message.Message{}, err, conn)

			return
		}

		if msg != nil {
			s.writeMsg(*msg, nil, conn)
		}
	}
}

// ProcessRequest - process request from client
func (s *Server) processRequest(_ context.Context, rawMsg string, clientInfo string) (*message.Message, error) {
	msg, err := message.ParseMessage(rawMsg)
	if err != nil {
		return nil, err
	}

	// switch by command from message
	switch msg.Command {
	case message.CommandError:
		return nil, errors.New("client requests to close connection")
	case message.CommandRequestChallenge:
		s.logger.Debug().Str("client", clientInfo).Msg("client requests challenge")

		hc, err := hashcash.New(s.config.HashBits, clientInfo)
		if err != nil {
			s.logger.Error().Err(err).Msg("err create hashcash")
			return nil, err
		}

		err = s.cacheRepo.SetWithTTL(hc.Key(), "", s.config.HashTTL)
		if err != nil {
			s.logger.Error().Err(err).Msg("err add rand to cache")
			return nil, err
		}

		msg := message.Message{
			Command: message.CommandResponseChallenge,
			Payload: string(hc.Header()),
		}

		return &msg, nil
	case message.CommandRequestResource:
		s.logger.Debug().Str("client", clientInfo).Any("payload", msg.Payload).Msg("client requests resource")

		// parse client's solution
		hc, err := hashcash.ParseHeader(msg.Payload)
		if err != nil {
			s.logger.Error().Err(err).Msg("parse payload")
			return nil, err
		}

		if !hc.EqualResource(clientInfo) {
			err := errors.New("resource is not equal")
			s.logger.Error().Err(err).Msg("validate error")
			return nil, err
		}

		if _, ok := s.cacheRepo.Get(hc.Key()); !ok {
			err := errors.New("hashcash not found")
			s.logger.Error().Err(err).Any("key", hc.Key()).Msg("validate error")
			return nil, err
		}

		if !hc.IsActual(s.config.HashTTL) {
			err := errors.New("date is not equal")
			s.logger.Error().Err(err).Msg("validate error")
			return nil, err
		}

		isHashCorrect, err := hc.Header().IsHashCorrect(hc.Bits())
		if err != nil {
			s.logger.Error().Err(err).Msg("validate error")
			return nil, err
		}
		if !isHashCorrect {
			err := errors.New("hash is not correct")
			s.logger.Error().Err(err).Msg("validate error")
			return nil, err
		}

		resource, err := s.quotesRepo.GetQuote()
		if err != nil {
			s.logger.Error().Err(err).Msg("get quote error")
			return nil, err
		}

		msg := message.Message{
			Command: message.CommandResponseResource,
			Payload: resource,
		}

		// delete hashcash from cache
		_ = s.cacheRepo.Delete(hc.Key())
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func (s *Server) writeMsg(msg message.Message, err error, w io.Writer) {
	if err != nil { // if error, send error message
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			s.logger.Error().Err(errWrite).Msg("io write error")
		}
	}

	if _, err := w.Write(msg.Bytes()); err != nil {
		s.logger.Error().Err(err).Msg("io write error")
	}
}
