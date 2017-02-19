package connectors

import "github.com/fluffle/goirc/client"

// NewIRCMessage constructs a new instance of an IRCMessage.
func NewIRCMessage(line *client.Line, self string) *IRCMessage {
	return &IRCMessage{line, self}
}

// IRCMessage is Phabulous's representation of an IRCMessage. This
// implementation is mainly a wrapper over the message struct provided by the
// IRC client library.
type IRCMessage struct {
	line *client.Line
	self string
}

// GetChannel returns the channel this message was posted on.
func (m *IRCMessage) GetChannel() string {
	return m.line.Target()
}

// GetUserID gets the nickname of the user who created this message.
func (m *IRCMessage) GetUserID() string {
	return m.line.Nick
}

// GetContent gets the message content.
func (m *IRCMessage) GetContent() string {
	return m.line.Text()
}

// GetProviderName returns the name of the provider this message was delivered
// by. Examples: slack, irc, etc.
func (m *IRCMessage) GetProviderName() string {
	return "irc"
}

// IsIM returns true if the message is a direct message sent to the bot.
func (m *IRCMessage) IsIM() bool {
	return m.line.Target() == m.line.Nick
}

// IsSelf returns true if the message was posted by the bot.
func (m *IRCMessage) IsSelf() bool {
	return m.line.Nick == m.self
}
