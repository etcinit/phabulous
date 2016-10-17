package modules

import (
	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/nlopes/slack"
)

// Service is an interface to the bot.
type Service interface {
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
	Excuse(*slack.MessageEvent, error)
	MakeGonduit() (*gonduit.Conn, error)
	MakeRTM() *slack.RTM
	MakeSlack() *slack.Client
	GetModules() []Module
}

// A Module provides a set of commands.
type Module interface {
	GetName() string
	GetCommands() []Command
}

// A Command provides access to a certain action.
type Command interface {
	GetUsage() string
	GetDescription() string
	GetMatchers() []string
	GetMentionMatchers() []string
	GetIMMatchers() []string
	GetHandler() Handler
}

// A Handler handles messages.
type Handler func(Service, *slack.MessageEvent, []string)
