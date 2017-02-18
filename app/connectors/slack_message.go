package connectors

import "github.com/nlopes/slack"

func NewSlackMessage(
	ev *slack.MessageEvent,
	self string,
	imChannelIDs *map[string]bool,
) *SlackMessage {
	return &SlackMessage{ev, self, imChannelIDs}
}

type SlackMessage struct {
	event        *slack.MessageEvent
	self         string
	imChannelIDs *map[string]bool
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

func (m *SlackMessage) IsIM() bool {
	_, ok := (*m.imChannelIDs)[m.event.Channel]

	return ok
}

func (m *SlackMessage) IsSelf() bool {
	return m.event.User == m.self
}
