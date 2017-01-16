package connectors

import "github.com/fluffle/goirc/client"

func NewIRCMessage(line *client.Line, self string) *IRCMessage {
	return &IRCMessage{line, self}
}

type IRCMessage struct {
	line *client.Line
	self string
}

func (m *IRCMessage) GetChannel() string {
	return m.line.Target()
}

func (m *IRCMessage) GetUserId() string {
	return m.line.Nick
}

func (m *IRCMessage) GetContent() string {
	return m.line.Text()
}

func (m *IRCMessage) GetProviderName() string {
	return "irc"
}

func (m *IRCMessage) IsIM() bool {
	return m.line.Target() == m.line.Nick
}

func (m *IRCMessage) IsSelf() bool {
	return m.line.Nick == m.self
}
