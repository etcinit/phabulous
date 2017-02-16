package messages

import (
	"github.com/fluffle/goirc/client"
	"github.com/nlopes/slack"
)

type Message interface {
	GetChannel() string
	GetUserId() string
	GetContent() string
	GetProviderName() string
}

func NewSlackMessage(ev *slack.MessageEvent) *SlackMessage {
	return &SlackMessage{ev}
}

func NewIRCMessage(line *client.Line) *IRCMessage {
	return &IRCMessage{line}
}

type SlackMessage struct {
	event *slack.MessageEvent
}

func (m *SlackMessage) GetChannel() string {
	return m.event.Channel
}

func (m *SlackMessage) GetUserId() string {
	return m.event.User
}

func (m *SlackMessage) GetContent() string {
	return m.event.Text
}

func (m *SlackMessage) GetProviderName() string {
	return "slack"
}

type IRCMessage struct {
	line *client.Line
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
