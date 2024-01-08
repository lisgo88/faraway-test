package message

import (
	"errors"
	"fmt"
	"strings"
)

// Command - command type
type Command int8

const (
	CommandError             Command = iota // using when something went wrong by client or server
	CommandRequestChallenge                 // request challenge from server (client->server)
	CommandResponseChallenge                // send challenge to client (server->client)
	CommandRequestResource                  // request resource from server (client->server)
	CommandResponseResource                 // send resource to client (server->client)

	DelimiterMessage = '\n' // DelimiterMessage - split messages from each other
	DelimiterCommand = ':'  // DelimiterCommand - split command and payload in message
)

// Message - message with command and payload
type Message struct {
	Command Command
	Payload string
}

// String - format message to string
func (m Message) String() string {
	return fmt.Sprintf("%d%c%s%c", m.Command, DelimiterCommand, m.Payload, DelimiterMessage)
}

// Bytes - format message to bytes
func (m Message) Bytes() []byte {
	return []byte(m.String())
}

// ParseMessage - parse message from request string
// raw message has `command:payload` format where command should be 0-4
func ParseMessage(msg string) (m Message, err error) {
	msg = strings.TrimSpace(msg)

	if len(msg) < 2 {
		return m, errors.New("incorrect message format")
	}

	switch msg[:2] {
	case fmt.Sprintf("0%c", DelimiterCommand):
		m.Command = CommandError
	case fmt.Sprintf("1%c", DelimiterCommand):
		m.Command = CommandRequestChallenge
	case fmt.Sprintf("2%c", DelimiterCommand):
		m.Command = CommandResponseChallenge
	case fmt.Sprintf("3%c", DelimiterCommand):
		m.Command = CommandRequestResource
	case fmt.Sprintf("4%c", DelimiterCommand):
		m.Command = CommandResponseResource
	default:
		return m, errors.New("unknown command")
	}

	if len(msg) > 2 {
		m.Payload = strings.TrimSpace(msg[2:])
	}

	return
}
