package interfaces

import "github.com/etcinit/phabulous/app/messages"

// A Module provides a set of commands.
type Module interface {
	// GetName returns the name of the Module.
	GetName() string

	// GetCommands returns all the Commands provided by this Module.
	GetCommands() []Command
}

// A Command provides access to a certain action.
type Command interface {
	// GetUsage returns a template of how a command should be invoked in a
	// similar fashion as a CLI utility would.
	GetUsage() string

	// GetDescription returns a short description of what the command does.
	GetDescription() string

	// GetMatchers returns regular expressions that match the command on
	// regular channels.
	GetMatchers() []string

	// GetMentionMatchers returns regular expressions that match the command
	// on regular channels when the bot is mention especifically.
	GetMentionMatchers() []string

	// GetIMMatchers returns regular expressions that match the command on
	// direct messages to the bot.
	GetIMMatchers() []string

	// GetHandler returns the handler function to be executed when the command
	// is matched.
	GetHandler() Handler
}

// A Handler handles messages.
type Handler func(Bot, messages.Message, []string)
