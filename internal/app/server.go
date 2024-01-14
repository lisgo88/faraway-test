package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"

	"github.com/lisgo88/faraway-test/internal/config"
	"github.com/lisgo88/faraway-test/internal/pkg/hashcash"
	"github.com/lisgo88/faraway-test/internal/pkg/message"
	"github.com/lisgo88/faraway-test/internal/pkg/repository"
)

type Server struct {
	quotesRepo repository.Quotes

	config config.ServerConfig
	logger zerolog.Logger
}

func NewServer(qr repository.Quotes, cfg config.ServerConfig, log zerolog.Logger) *Server {
	return &Server{
		quotesRepo: qr,
		config:     cfg,
		logger:     log,
	}
}

// Run - start server
func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer func() {
		_ = listener.Close()
	}()

	s.logger.Info().Any("config", listener.Addr()).Msg("TCP server starting...")

	// graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		s.logger.Info().Msg("Received termination signal. Closing server...")

		cancel()
		_ = listener.Close()
	}()

	for {
		// listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				s.logger.Info().Msg("Server is shutting down.")
				return nil
			default:
				s.logger.Error().Err(err).Msg("accept connection error")
				continue
			}
		}

		// set timeout for connection
		_ = conn.SetReadDeadline(time.Now().Add(s.config.TimeOut * time.Second))

		// handle connections
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	ctx = context.WithValue(ctx, "clientID", conn.RemoteAddr().String())
	s.logger.Info().Any("clientID", ctx.Value("clientID")).Msg("new client")

	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Any("requestID", ctx.Value("requestID")).Msg("connection handling canceled")
			return
		default:
			ctx = context.WithValue(ctx, "requestID", uuid.NewV4())

			// create hashcash
			hc, err := hashcash.New(ctx, s.config.HashBits)
			if err != nil {
				s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("create hashcash")
				return
			}

			// send hashcash to client
			msg := &message.Message{
				RequestID: ctx.Value("requestID").(uuid.UUID).String(),
				ClientID:  ctx.Value("clientID").(string),
				Challenge: string(hc.Header()),
			}

			data, err := proto.Marshal(msg)
			if err != nil {
				s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("marshal msg")
				s.writeMsg(nil, err, conn)

				return
			}

			_, err = conn.Write(data)
			if err != nil {
				s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("write msg")
				s.writeMsg(nil, err, conn)

				return
			}

			length, _ := conn.Read(buffer)
			clientMsg := &message.Message{}
			if err := proto.Unmarshal(buffer[:length], clientMsg); err != nil {
				s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("unmarshal msg")
				s.writeMsg(nil, err, conn)

				return
			}

			clientHc, err := hashcash.ParseHeader(clientMsg.Challenge)
			if err != nil {
				s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("parse hashcash")
				s.writeMsg(nil, err, conn)

				return
			}

			if ok, err := clientHc.IsValid(ctx, s.config.HashTTL); ok && err == nil {
				quot, _ := s.quotesRepo.GetQuote()
				msg.Payload = quot

				data, err := proto.Marshal(msg)
				if err != nil {
					s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("marshal msg")
					s.writeMsg(nil, err, conn)

					return
				}

				s.writeMsg(data, nil, conn)
			} else {
				s.logger.Error().Err(err).Any("requestID", ctx.Value("requestID")).Msg("valid hashcash")
				s.writeMsg(nil, err, conn)

				return
			}
		}
	}
}

func (s *Server) writeMsg(msg []byte, err error, w io.Writer) {
	if err != nil { // if error, send error message
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			s.logger.Error().Err(errWrite).Msg("io write error")
		}
	}

	if _, err := w.Write(msg); err != nil {
		s.logger.Error().Err(err).Msg("io write error")
	}
}
