package slacker

import "github.com/nlopes/slack"

// SlackService provides access to the Slack service.
type SlackService struct {
	Slack *slack.Slack
}
