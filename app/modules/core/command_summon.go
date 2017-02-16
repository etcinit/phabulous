package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
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
func (c *SummonCommand) GetHandler() modules.Handler {
	return func(s modules.Service, ev *slack.MessageEvent, matches []string) {
		s.StartTyping(ev.Channel)

		if len(matches) < 2 {
			return
		}

		conn, err := s.MakeGonduit()
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		res, err := conn.DifferentialQuery(requests.DifferentialQueryRequest{
			IDs: []uint64{uint64(id)},
		})
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		if len(*res) == 0 {
			s.Post(
				ev.Channel,
				"Revision not found.",
				messages.IconDefault,
				true,
			)

			return
		}

		if len((*res)[0].Reviewers) == 0 {
			s.Post(
				ev.Channel,
				"Revision has no reviewers.",
				messages.IconDefault,
				true,
			)

			return
		}

		reviewerNames := []string{}

		for _, reviewerPHID := range (*res)[0].Reviewers {
			nameRes, err := conn.PHIDQuerySingle(reviewerPHID)
			if err != nil {
				s.Excuse(ev, err)
				return
			}

			reviewerNames = append(reviewerNames, "@"+(*nameRes).Name)
		}

		userInfo, err := s.MakeRTM().GetUserInfo(ev.User)
		if err != nil {
			s.Excuse(ev, err)
			return
		}

		s.Post(
			ev.Channel,
			fmt.Sprintf(
				"*@%s summons %s to review D%s:*\n_%s (%s)_\n%s",
				userInfo.Name,
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
