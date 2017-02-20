package connectors

import (
	"strings"

	"github.com/etcinit/phabulous/app/messages"
)

// Post posts a text message.
func (c *IRCConnector) Post(
	channelName string,
	storyText string,
	icon messages.Icon,
	asUser bool,
) {
	for _, line := range strings.Split(storyText, "\n") {
		c.client.Privmsgln(channelName, line)
	}
}

// PostImage posts a message with an attached image.
func (c *IRCConnector) PostImage(
	channelName string,
	storyText string,
	imageURL string,
	icon messages.Icon,
	asUser bool,
) {
	for _, line := range strings.Split(storyText, "\n") {
		c.client.Privmsgln(channelName, line)
	}

	c.client.Privmsgln(channelName, imageURL)
}

// PostOnFeed posts a message on the bot's "feed" channel.
func (c *IRCConnector) PostOnFeed(storyText string) error {
	if c.GetFeedChannel() == "" {
		return ErrMissingFeedChannel
	}

	for _, line := range strings.Split(storyText, "\n") {
		c.Post(c.GetFeedChannel(), line, messages.IconDefault, false)
	}

	return nil
}

// GetFeedChannel returns the default channel for the bot.
func (c *IRCConnector) GetFeedChannel() string {
	return c.config.GetString("channels.feed")
}
