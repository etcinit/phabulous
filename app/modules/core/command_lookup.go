package core

import (
	"fmt"

	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/nlopes/slack"
)

// LookupCommand allows users to lookup objects from Phabricator.
type LookupCommand struct{}

// GetUsage returns the usage of this command.
func (c *LookupCommand) GetUsage() string {
	return "lookup (Txxx|Dxxx)"
}

// GetDescription returns the description of this command.
func (c *LookupCommand) GetDescription() string {
	return "Looks up a task or revision by its number."
}

// GetMatchers returns the matchers for this command.
func (c *LookupCommand) GetMatchers() []string {
	return []string{
		"^([T|D][0-9]{1,16})$",
	}
}

// GetIMMatchers returns IM matchers for this command.
func (c *LookupCommand) GetIMMatchers() []string {
	return []string{
		"^lookup\\s+([T|D][0-9]{1,16})$",
		"^([T|D][0-9]{1,16})$",
	}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (c *LookupCommand) GetMentionMatchers() []string {
	return []string{
		"lookup\\s+([T|D][0-9]{1,16})",
		"([T|D][0-9]{1,16})",
	}
}

// GetHandler returns the handler for this command.
func (c *LookupCommand) GetHandler() modules.Handler {
	return func(s modules.Service, ev *slack.MessageEvent, matches []string) {
		s.StartTyping(ev.Channel)

		conn, err := s.MakeGonduit()
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		res, err := conn.PHIDLookupSingle(matches[1])
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		if res == nil {
			s.Post(
				ev.Channel,
				fmt.Sprintf("I couldn't find %s", matches[1]),
				messages.IconDefault,
				true,
			)
			return
		}

		s.Post(
			ev.Channel,
			fmt.Sprintf("*%s* (%s): %s", res.FullName, res.Status, res.URI),
			messages.IconTasks,
			true,
		)
	}
}
