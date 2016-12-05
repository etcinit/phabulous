package connectors

import (
	"errors"

	"github.com/etcinit/phabulous/app/messages"
	"github.com/nlopes/slack"
)

var (
	// ErrMissingFeedChannel is used when the feed channel is not configured.
	ErrMissingFeedChannel = errors.New("Missing feed channel")
)

// Post posts a simple message to Slack. Most parameters are set to
// defaults.
func (c *SlackConnector) Post(
	channelName string,
	storyText string,
	icon messages.Icon,
	asUser bool,
) {
	user := c.config.GetString("slack.username")

	if c.slackInfo != nil {
		user = c.slackInfo.User.Name
	}

	if c.config.GetBool("slack.as-user") == true {
		asUser = true
	}

	c.slack.PostMessage(
		channelName,
		storyText,
		slack.PostMessageParameters{
			Username: user,
			IconURL:  string(icon),
			AsUser:   asUser,
		},
	)
}

// PostImage posts a simple message to Slack, with an image. Most parameters
// are set to defaults.
func (c *SlackConnector) PostImage(
	channelName string,
	storyText string,
	imageURL string,
	icon messages.Icon,
	asUser bool,
) {
	user := c.config.GetString("slack.username")

	if c.slackInfo != nil {
		user = c.slackInfo.User.Name
	}

	if c.config.GetBool("slack.as-user") == true {
		asUser = true
	}

	c.slack.PostMessage(
		channelName,
		storyText,
		slack.PostMessageParameters{
			Username: user,
			IconURL:  string(icon),
			AsUser:   asUser,
			Attachments: []slack.Attachment{
				{ImageURL: imageURL},
			},
		},
	)
}

// PostOnFeed posts a message to Slack on the default bot channel.
func (c *SlackConnector) PostOnFeed(storyText string) error {
	if c.GetFeedChannel() == "" {
		return ErrMissingFeedChannel
	}

	c.Post(c.GetFeedChannel(), storyText, messages.IconDefault, false)

	return nil
}

// GetFeedChannel returns the default channel for the bot.
func (c *SlackConnector) GetFeedChannel() string {
	return c.config.GetString("channels.feed")
}
