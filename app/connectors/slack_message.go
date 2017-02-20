package connectors

import "github.com/nlopes/slack"

// NewSlackMessage constructs a new instance of an SlackMessage.
func NewSlackMessage(
	ev *slack.MessageEvent,
	self string,
	imChannelIDs *map[string]bool,
) *SlackMessage {
	return &SlackMessage{ev, self, imChannelIDs}
}

// SlackMessage is Phabulous's representation of a Slack message. This
// implementation is mainly a wrapper over the message struct provided by the
// Slack client library, along with some additional metadata.
type SlackMessage struct {
	event        *slack.MessageEvent
	self         string
	imChannelIDs *map[string]bool
}

// GetChannel returns the channel this message was posted on.
func (m *SlackMessage) GetChannel() string {
	return m.event.Channel
}

// GetUserID gets the Slack account ID of the user who created this message.
func (m *SlackMessage) GetUserID() string {
	return m.event.User
}

// GetContent gets the message content.
func (m *SlackMessage) GetContent() string {
	return m.event.Text
}

// GetProviderName returns the name of the provider this message was delivered
// by. Examples: slack, irc, etc.
func (m *SlackMessage) GetProviderName() string {
	return "slack"
}

// IsIM returns true if the message is a direct message sent to the bot.
func (m *SlackMessage) IsIM() bool {
	_, ok := (*m.imChannelIDs)[m.event.Channel]

	return ok
}

// IsSelf returns true if the message was posted by the bot.
func (m *SlackMessage) IsSelf() bool {
	return m.event.User == m.self
}
