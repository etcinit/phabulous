package core

import (
	"fmt"

	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/nlopes/slack"
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
func (c *ModulesCommand) GetHandler() modules.Handler {
	return func(s modules.Service, ev *slack.MessageEvent, matches []string) {
		message := "Loaded modules:\n"

		for _, module := range s.GetModules() {

			message = message + fmt.Sprintf(
				"• *%s*\n",
				module.GetName(),
			)
		}

		s.Post(ev.Channel, message, messages.IconTasks, true)
	}
}
