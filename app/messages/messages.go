package messages

import "github.com/nlopes/slack"

type Message interface {
	GetChannel() string
	GetUserId() string
	GetContent() string
	GetProviderName() string
}

func NewSlackMessage(ev *slack.MessageEvent) *SlackMessage {
	return &SlackMessage{ev}
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
