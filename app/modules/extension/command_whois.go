package extension

import (
	"fmt"

	"github.com/etcinit/gonduit"
	gonduitRequests "github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/gonduit/extensions"
	"github.com/etcinit/phabulous/app/gonduit/extensions/requests"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/nlopes/slack"
)

// WhoisCommand allows one to send test messages to the feed channel.
type WhoisCommand struct{}

// GetUsage returns the usage of this command.
func (t *WhoisCommand) GetUsage() string {
	return "whois <slack|phabricator> <username>"
}

// GetDescription returns the description of this command.
func (t *WhoisCommand) GetDescription() string {
	return "Gets the name of a Slack user in Phabricator and vice-versa."
}

// GetMatchers returns the matchers for this command.
func (t *WhoisCommand) GetMatchers() []string {
	return []string{}
}

// GetIMMatchers returns IM matchers for this command.
func (t *WhoisCommand) GetIMMatchers() []string {
	return []string{"whois (slack|phabricator) ([0-9a-zA-Z]+)"}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (t *WhoisCommand) GetMentionMatchers() []string {
	return []string{"whois (slack|phabricator) ([0-9a-zA-Z]+)"}
}

// GetHandler returns the handler for this command.
func (t *WhoisCommand) GetHandler() modules.Handler {
	return func(s modules.Service, ev *slack.MessageEvent, matches []string) {
		s.StartTyping(ev.Channel)

		conn, err := s.MakeGonduit()

		if err != nil {
			s.Excuse(ev, err)
			return
		}

		if matches[1] == "slack" {
			fromSlack(s, conn, ev, matches[2])
		} else if matches[1] == "phabricator" {
			toSlack(s, conn, ev, matches[2])
		}
	}
}

func toSlack(
	s modules.Service,
	conn *gonduit.Conn,
	ev *slack.MessageEvent,
	username string,
) {
	users, err := conn.UserQuery(gonduitRequests.UserQueryRequest{
		Usernames: []string{username},
	})
	if err != nil {
		s.Excuse(ev, err)
		return
	}

	if len(*users) == 0 {
		s.Post(
			ev.Channel,
			"I was unable to find a user with that name on Phabricator.",
			messages.IconTasks,
			true,
		)
		return
	}

	accounts, err := extensions.PhabulousToSlack(
		conn,
		requests.PhabulousToSlackRequest{
			UserPHIDs: []string{(*users)[0].PHID},
		},
	)
	if err != nil {
		s.Excuse(ev, err)
		return
	}

	if len(*accounts) == 0 {
		s.Post(
			ev.Channel,
			"I was unable to find a Phabricator user linked with that "+
				"Slack account. Make sure they are linked under "+
				"_External Accounts_ in the user's Phabricator settings.",
			messages.IconTasks,
			true,
		)
		return
	}

	client := s.MakeSlack()
	slackUsers, err := client.GetUsers()
	if err != nil {
		s.Excuse(ev, err)
		return
	}

	var foundUser *slack.User
	for _, user := range slackUsers {
		if user.Name == username {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		s.Post(
			ev.Channel,
			"I was unable to find a user with that name on this Slack "+
				"organization",
			messages.IconTasks,
			true,
		)
		return
	}

	s.Post(
		ev.Channel,
		fmt.Sprintf(
			"*%s* is known as *%s* on Slack.",
			username,
			foundUser.Name,
		),
		messages.IconTasks,
		true,
	)
}

func fromSlack(
	s modules.Service,
	conn *gonduit.Conn,
	ev *slack.MessageEvent,
	username string,
) {
	client := s.MakeSlack()

	users, err := client.GetUsers()
	if err != nil {
		s.Excuse(ev, err)
		return
	}

	var foundUser *slack.User
	for _, user := range users {
		if user.Name == username {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		s.Post(
			ev.Channel,
			"I was unable to find a user with that name on this Slack "+
				"organization",
			messages.IconTasks,
			true,
		)
		return
	}

	res, err := extensions.PhabulousFromSlack(
		conn,
		requests.PhabulousFromSlackRequest{
			AccountIDs: []string{foundUser.ID},
		},
	)
	if err != nil {
		s.Excuse(ev, err)
		return
	}

	if len(*res) == 0 {
		s.Post(
			ev.Channel,
			"I was unable to find a Phabricator user linked with that "+
				"Slack account. Make sure they are linked under "+
				"_External Accounts_ in the user's Phabricator settings.",
			messages.IconTasks,
			true,
		)
		return
	}

	userPHIDs := []string{}
	for _, userInfo := range *res {
		userPHIDs = append(userPHIDs, userInfo.UserPHID)
	}

	res2, err := conn.UserQuery(gonduitRequests.UserQueryRequest{
		PHIDs: userPHIDs,
	})

	if err != nil {
		s.Excuse(ev, err)
		return
	}

	for _, user := range *res2 {
		s.Post(
			ev.Channel,
			fmt.Sprintf(
				"*%s* is known as *%s* on Phabricator.",
				username,
				user.UserName,
			),
			messages.IconTasks,
			true,
		)
	}
}
