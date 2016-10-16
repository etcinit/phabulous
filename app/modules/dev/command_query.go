package dev

import (
	"fmt"

	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/nlopes/slack"
)

// QueryCommand allows one to send test messages to the feed channel.
type QueryCommand struct{}

// GetUsage returns the usage of this command.
func (t *QueryCommand) GetUsage() string {
	return "dev:query"
}

// GetDescription returns the description of this command.
func (t *QueryCommand) GetDescription() string {
	return "Lists all available conduit endpoints."
}

// GetMatchers returns the matchers for this command.
func (t *QueryCommand) GetMatchers() []string {
	return []string{"dev:query"}
}

// GetIMMatchers returns IM matchers for this command.
func (t *QueryCommand) GetIMMatchers() []string {
	return []string{}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (t *QueryCommand) GetMentionMatchers() []string {
	return []string{}
}

// GetHandler returns the handler for this command.
func (t *QueryCommand) GetHandler() modules.Handler {
	return func(s modules.Service, ev *slack.MessageEvent, matches []string) {
		conn, err := s.MakeGonduit()

		if err != nil {
			s.Excuse(ev, err)
			return
		}

		res, err := conn.ConduitQuery()
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		message := "Available Conduit methods:\n"

		for methodName, method := range *res {
			message = message + fmt.Sprintf(
				"• *%s*:\n\t_%s_\n",
				methodName,
				method.Description,
			)
		}

		s.Post(ev.Channel, message, messages.IconTasks, true)
	}
}
