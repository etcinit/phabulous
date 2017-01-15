package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/gonduit/extensions"
	phabulousRequests "github.com/etcinit/phabulous/app/gonduit/extensions/requests"
	"github.com/etcinit/phabulous/app/gonduit/extensions/responses"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/nlopes/slack"
)

// SummonCommand allows users to summon reviewers.
type SummonCommand struct{}

// GetUsage returns the usage of this command.
func (c *SummonCommand) GetUsage() string {
	return "summon Dxxx"
}

// GetDescription returns the description of this command.
func (c *SummonCommand) GetDescription() string {
	return "Asks reviewers of a revision to review it."
}

// GetMatchers returns the matchers for this command.
func (c *SummonCommand) GetMatchers() []string {
	return []string{}
}

// GetIMMatchers returns IM matchers for this command.
func (c *SummonCommand) GetIMMatchers() []string {
	return []string{
		"summon\\s+D([0-9]{1,16})",
	}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (c *SummonCommand) GetMentionMatchers() []string {
	return []string{
		"summon\\s+D([0-9]{1,16})",
	}
}

// GetHandler returns the handler for this command.
func (c *SummonCommand) GetHandler() interfaces.Handler {
	return func(s interfaces.Bot, m messages.Message, matches []string) {
		s.StartTyping(m.GetChannel())

		if len(matches) < 2 {
			return
		}

		conn, err := s.GetGonduit()
		if err != nil {
			s.Excuse(m, err)
			return
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			s.Excuse(m, err)
			return
		}

		res, err := conn.DifferentialQuery(requests.DifferentialQueryRequest{
			IDs: []uint64{uint64(id)},
		})
		if err != nil {
			s.Excuse(m, err)
			return
		}

		if len(*res) == 0 {
			s.Post(
				m.GetChannel(),
				"Revision not found.",
				messages.IconDefault,
				true,
			)

			return
		}

		if len((*res)[0].Reviewers) == 0 {
			s.Post(
				m.GetChannel(),
				"Revision has no reviewers.",
				messages.IconDefault,
				true,
			)

			return
		}

		var slackMap *responses.PhabulousToSlackResponse
		var slackUsers []slack.User

		if sb, ok := s.(interfaces.SlackBot); ok {
			slackMap, err = extensions.PhabulousToSlack(
				conn,
				phabulousRequests.PhabulousToSlackRequest{
					UserPHIDs: (*res)[0].Reviewers,
				},
			)
			if err != nil {
				s.Excuse(m, err)
				return
			}

			slackUsers, err = sb.GetSlack().GetUsers()
			if err != nil {
				s.Excuse(m, err)
				return
			}
		}

		reviewerNames := []string{}

		for _, reviewerPHID := range (*res)[0].Reviewers {
			if _, ok := s.(interfaces.SlackBot); ok {
				if slackUserInfo, ok := (*slackMap)[reviewerPHID]; ok {
					var foundUser *slack.User
					for _, user := range slackUsers {
						if user.ID == slackUserInfo.AccountID {
							foundUser = &user
							break
						}
					}

					if foundUser != nil {
						reviewerNames = append(
							reviewerNames,
							fmt.Sprintf(
								"@%s :slack:",
								foundUser.Name,
							),
						)

						continue
					}
				}
			}

			nameRes, err := conn.PHIDQuerySingle(reviewerPHID)
			if err != nil {
				s.Excuse(m, err)
				return
			}

			reviewerNames = append(reviewerNames, "@"+(*nameRes).Name)
		}

		userName, err := s.GetUsername(m.GetUserId())
		if err != nil {
			s.Excuse(m, err)
			return
		}

		s.Post(
			m.GetChannel(),
			fmt.Sprintf(
				"*@%s summons %s to review D%s:*\n_%s (%s)_\n%s",
				userName,
				strings.Join(reviewerNames, ", "),
				matches[1],
				(*res)[0].Title,
				(*res)[0].StatusName,
				(*res)[0].URI,
			),
			messages.IconDefault,
			true,
		)
	}
}
