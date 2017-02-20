package interfaces

import (
	"github.com/etcinit/gonduit"
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
	Poster

	// GetUsername returns the username of a user in the current network when
	// provided with a user ID.
	//
	// This is mainly useful if the chat network has a notion of user IDs that
	// uniquely identify a user regardless of username changes (i.e. Slack).
	//
	GetUsername(userID string) (string, error)

	// StartTyping will cause the bot to show a typing indicator
	// ("X is typing...") if the network supports it. Otherwise, it will simply
	// ignore it.
	StartTyping(channelID string)

	// GetUsageHandler returns a handler to be used for when no other handlers are
	// matched. This handler usually posts some for of help message.
	GetUsageHandler() Handler

	// Excuse can be used as an error reporter by commands. It posts to the
	// channel a message was received from that an error occurred and logs the
	// error using the application logger.
	Excuse(Message, error)

	// GetGonduit returns an instance of the Conduit client.
	GetGonduit() (*gonduit.Conn, error)

	// GetConfig returns an instance of the configuration object.
	GetConfig() *confer.Config

	// GetModules returns a slice of all the modules loaded by this Bot.
	GetModules() []Module

	// GetHandlers returns the currently active handlers on the connector.
	GetHandlers() []HandlerTuple

	// GetIMHandlers returns the currently active IM handlers on the connector.
	GetIMHandlers() []HandlerTuple
}

// A SlackBot is just like a Bot, but it also provides access to the Slack API.
// This might be needed by some commands that rely on Slack-specific
// functionality.
type SlackBot interface {
	Bot

	// GetSlack returns an instance of the Slack client for the bot's network.
	GetSlack() *slack.Client
}
