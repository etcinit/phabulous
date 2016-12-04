package core

import (
	"fmt"
	"strings"

	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
)

// HelpCommand allows one to send test messages to the feed channel.
type HelpCommand struct{}

// GetUsage returns the usage of this command.
func (c *HelpCommand) GetUsage() string {
	return "help"
}

// GetDescription returns the description of this command.
func (c *HelpCommand) GetDescription() string {
	return "Display a listing of all available commands."
}

// GetMatchers returns the matchers for this command.
func (c *HelpCommand) GetMatchers() []string {
	return []string{}
}

// GetIMMatchers returns IM matchers for this command.
func (c *HelpCommand) GetIMMatchers() []string {
	return []string{"^help$"}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (c *HelpCommand) GetMentionMatchers() []string {
	return []string{"help"}
}

// GetHandler returns the handler for this command.
func (c *HelpCommand) GetHandler() modules.Handler {
	return func(s modules.Service, m messages.Message, matches []string) {
		message := "Available commands:\n"

		for _, module := range s.GetModules() {
			for _, command := range module.GetCommands() {
				methods := []string{}

				if len(command.GetMatchers()) > 0 {
					methods = append(methods, "Passive")
				}

				if len(command.GetIMMatchers()) > 0 {
					methods = append(methods, "DM")
				}

				if len(command.GetMentionMatchers()) > 0 {
					methods = append(methods, "Mention")
				}

				message = message + fmt.Sprintf(
					"• *%s*\t[%s]\n\t_%s_\n",
					command.GetUsage(),
					strings.Join(methods, ", "),
					command.GetDescription(),
				)
			}
		}

		s.Post(m.GetChannel(), message, messages.IconTasks, true)
	}
}
