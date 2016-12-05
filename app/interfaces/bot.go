package interfaces

import (
	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// A Bot provides most methods and services needed by command handlers to
// perform their action.
//
// The interface is high-level to allow for implementations on different
// networks and services.
//
type Bot interface {
	Post(string, string, messages.Icon, bool)
	PostImage(
		channelName string,
		storyText string,
		imageURL string,
		icon messages.Icon,
		asUser bool,
	)
	PostOnFeed(string)
	StartTyping(string)
	GetUsername(string) (string, error)
	Excuse(messages.Message, error)
	GetGonduit() (*gonduit.Conn, error)
	GetConfig() *confer.Config
	GetModules() []Module
}

// A SlackBot is just like a Bot, but it also provides access to the Slack API.
// This might be needed by some commands that rely on Slack-specific
// functionality.
type SlackBot interface {
	Bot
	GetSlack() *slack.Client
}
