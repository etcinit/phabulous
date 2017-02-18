package core

import (
	"fmt"

	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
)

// ModulesCommand allows one to send test messages to the feed channel.
type ModulesCommand struct{}

// GetUsage returns the usage of this command.
func (c *ModulesCommand) GetUsage() string {
	return "modules"
}

// GetDescription returns the description of this command.
func (c *ModulesCommand) GetDescription() string {
	return "Display a listing of all modules loaded."
}

// GetMatchers returns the matchers for this command.
func (c *ModulesCommand) GetMatchers() []string {
	return []string{}
}

// GetIMMatchers returns IM matchers for this command.
func (c *ModulesCommand) GetIMMatchers() []string {
	return []string{"^modules$"}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (c *ModulesCommand) GetMentionMatchers() []string {
	return []string{"modules"}
}

// GetHandler returns the handler for this command.
func (c *ModulesCommand) GetHandler() interfaces.Handler {
	return func(s interfaces.Bot, m interfaces.Message, matches []string) {
		message := "Loaded modules:\n"

		for _, module := range s.GetModules() {

			message = message + fmt.Sprintf(
				"• *%s*\n",
				module.GetName(),
			)
		}

		s.Post(m.GetChannel(), message, messages.IconTasks, true)
	}
}
