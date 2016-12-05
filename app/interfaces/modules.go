package interfaces

import "github.com/etcinit/phabulous/app/messages"

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
type Handler func(Bot, messages.Message, []string)
