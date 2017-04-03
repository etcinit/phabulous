package core

import (
	"bytes"
	"fmt"

	"github.com/etcinit/gonduit/constants"
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
)

// TaskCommand lists your tasks
type TaskCommand struct{}

// GetUsage .
func (c *TaskCommand) GetUsage() string {
	return "listtasks <username>"
}

// GetDescription .
func (c *TaskCommand) GetDescription() string {
	return "Lists currently-open tasks for a given Phabricator user. If no username is provided, Slack username is used."
}

// GetMatchers .
func (c *TaskCommand) GetMatchers() []string {
	return []string{
		"^listtasks\\s*(.*)$",
	}
}

// GetIMMatchers .
func (c *TaskCommand) GetIMMatchers() []string {
	return []string{
		"^listtasks\\s*(.*)$",
	}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (c *TaskCommand) GetMentionMatchers() []string {
	return []string{
		"^listtasks\\s*(.*)$",
	}
}

// GetHandler -- returns a handler that mostly ignores what was
// matched and just returns the user's list of open tasks.
func (c *TaskCommand) GetHandler() interfaces.Handler {
	return func(s interfaces.Bot, m interfaces.Message, matches []string) {
		s.StartTyping(m.GetChannel())

		conn, err := s.GetGonduit()
		if err != nil {
			s.Excuse(m, err)
			return
		}

		username := matches[1]
		if len(username) == 0 {
			username, err = s.GetUsername(m.GetUserID())
			if err != nil || username == "" {
				s.Excuse(m, err)
				return
			}
		}
		nameres, err := conn.UserQuery(
			requests.UserQueryRequest{Usernames: []string{username}},
		)
		if err != nil || nameres == nil || len(*nameres) == 0 {
			message := fmt.Sprintf(`
Couldn't find Phabricator user for Slack username *@%s*. Typically
this means your Slack username doesn't match your Phabricator
username.
			`,
				username,
			)
			s.Post(
				m.GetChannel(),
				message,
				messages.IconDefault,
				true,
			)
			return
		}
		ownerPHID := (*nameres)[0].PHID

		res, err := conn.ManiphestQuery(requests.ManiphestQueryRequest{
			OwnerPHIDs: []string{ownerPHID},
			Status:     constants.ManiphestTaskStatusOpen,
			Order:      constants.ManiphestQueryOrderPriority,
		})

		if err != nil {
			s.Excuse(m, err)
			return
		}

		if res == nil {
			s.Post(
				m.GetChannel(),
				fmt.Sprintf("<@%s> doesn't appear to have any open tasks", username),
				messages.IconDefault,
				true,
			)
			return
		}

		var buffer bytes.Buffer
		for _, value := range *res {
			buffer.WriteString(
				fmt.Sprintf("*<%s|T%s>* -- %v\n", value.URI, value.ID, value.Title),
			)
		}
		s.Post(
			m.GetChannel(),
			buffer.String(),
			messages.IconTasks,
			true,
		)
	}
}
