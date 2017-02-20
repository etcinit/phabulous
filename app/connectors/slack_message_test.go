package connectors

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func Test_NewSlackMessage(t *testing.T) {
	event := &slack.MessageEvent{}

	NewSlackMessage(event, "someid", &map[string]bool{})
}

func Test_SlackMessage_GetChannel(t *testing.T) {
	event := &slack.MessageEvent{}
	event.Channel = "somechannel"

	message := NewSlackMessage(event, "someid", &map[string]bool{})

	assert.Equal(t, "somechannel", message.GetChannel())
}

func Test_SlackMessage_GetUserID(t *testing.T) {
	event := &slack.MessageEvent{}
	event.User = "someuserid"

	message := NewSlackMessage(event, "someid", &map[string]bool{})

	assert.Equal(t, "someuserid", message.GetUserID())
}

func Test_SlackMessage_GetContent(t *testing.T) {
	event := &slack.MessageEvent{}
	event.Text = "This is art"

	message := NewSlackMessage(event, "someid", &map[string]bool{})

	assert.Equal(t, "This is art", message.GetContent())
}

func Test_SlackMessage_GetProviderName(t *testing.T) {
	event := &slack.MessageEvent{}
	message := NewSlackMessage(event, "someid", &map[string]bool{})

	assert.Equal(t, "slack", message.GetProviderName())
}

func Test_SlackMessage_IsIM(t *testing.T) {
	event := &slack.MessageEvent{}
	event.Channel = "notifications"

	message := NewSlackMessage(event, "someid", &map[string]bool{})

	assert.False(t, message.IsIM())

	(*message.imChannelIDs)["notifications"] = true

	assert.True(t, message.IsIM())
}

func Test_SlackMessage_IsSelf(t *testing.T) {
	event := &slack.MessageEvent{}
	event.User = "someid"

	message := NewSlackMessage(event, "someid", &map[string]bool{})

	assert.True(t, message.IsSelf())
}
