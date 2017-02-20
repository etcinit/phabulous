package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/gonduit/entities"
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/connectors/utilities"
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
	return func(s interfaces.Bot, m interfaces.Message, matches []string) {
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

		revision := (*res)[0]

		if len(revision.Reviewers) == 0 {
			s.Post(
				m.GetChannel(),
				"Revision has no reviewers.",
				messages.IconDefault,
				true,
			)

			return
		}

		reviewerNames, err := c.getReviewerNames(s, conn, revision)
		if err != nil {
			s.Excuse(m, err)
			return
		}

		userName, err := s.GetUsername(m.GetUserID())
		if err != nil {
			s.Excuse(m, err)
			return
		}

		s.Post(
			m.GetChannel(),
			fmt.Sprintf(
				"*@%s summons %s to review D%s:*\n_%s (%s)_\n%s",
				userName,
				c.namesToString(reviewerNames),
				matches[1],
				revision.Title,
				revision.StatusName,
				revision.URI,
			),
			messages.IconDefault,
			true,
		)
	}
}

// namesToString converts the list of names into a single string, providing a
// default value for when the list is empty.
func (c *SummonCommand) namesToString(names []string) string {
	if len(names) == 0 {
		return "(nobody)"
	}

	return strings.Join(names, ", ")
}

// getReviewerNames does a lot of magic to give us a pretty list of reviewers
// that we can return in a message. Phbricator usernames will be used, projects
// will be expanded, and Slack usernames will be resolved if possible.
func (c *SummonCommand) getReviewerNames(
	bot interfaces.Bot,
	conn *gonduit.Conn,
	revision *entities.DifferentialRevision,
) ([]string, error) {
	slackMap, slackUsers, err := c.lookupSlackMap(bot, conn, revision)

	if err != nil {
		return nil, err
	}

	reviewerNames := []string{}
	reviewerMap, err := c.getReviewerPHIDs(bot, conn, revision)

	if err != nil {
		return nil, err
	}

	includeSelf := bot.GetConfig().GetBool("core.summon.includeSelf")

	for reviewerPHID, reviewerName := range reviewerMap {
		// Prevent the author from embarassing themselves.
		if !includeSelf && revision.AuthorPHID == reviewerPHID {
			continue
		}

		// Attempt to match on Slack.
		if m, ok := c.findOnSlack(bot, slackMap, &slackUsers, reviewerPHID); ok {
			reviewerNames = append(reviewerNames, m)

			continue
		}

		// Otherwise, use their Phabricator username.
		reviewerNames = append(reviewerNames, "@"+reviewerName)
	}

	return reviewerNames, nil
}

// getReviewerPHIDs gets a mapping of PHIDs and usernames. Projects are expanded
// if needed and configured.
//
// We do some batch requests, which save time and resources, but they might be
// limited by pagination. If the diff has many attached users and projects,
// the bot might produce partial results.
//
// The current implementation assumes diffs won't have a ridiculous amount of
// reviewers attached to them, and that if projects are used, they have a
// moderate amount of users as well.
func (c *SummonCommand) getReviewerPHIDs(
	bot interfaces.Bot,
	conn *gonduit.Conn,
	revision *entities.DifferentialRevision,
) (map[string]string, error) {
	reviewerMap := map[string]string{}
	projects := []string{}
	expandProjects := bot.GetConfig().GetBool("core.summon.expandProjects")

	// We query all PHIDs in batch to avoid spamming the Phabricator server with
	// individual requests.
	allRes, err := conn.PHIDQuery(requests.PHIDQueryRequest{
		PHIDs: revision.Reviewers,
	})

	if err != nil {
		return nil, err
	}

	for reviewerPHID, queryResult := range allRes {
		// If the PHID is a project, we will keep it for later.
		if queryResult.Type == "PROJ" && expandProjects {
			projects = append(projects, reviewerPHID)

			continue
		}

		reviewerMap[reviewerPHID] = queryResult.Name
	}

	// If any of the reviewers was a project, we will do some additional
	// processing.
	if len(projects) > 0 {
		projectMembers := map[string]bool{}
		projectMembersList := []string{}

		// Batch request all projects.
		projRes, err := conn.ProjectQuery(requests.ProjectQueryRequest{
			PHIDs: projects,
		})

		if err != nil {
			return nil, err
		}

		// Obtain a Set union of the members of all projects.
		for _, project := range projRes.Data {
			for _, memberPHID := range project.Members {
				projectMembers[memberPHID] = true
			}
		}

		// Extract all Set keys (PHIDs) into a list.
		for projectMember := range projectMembers {
			// Small optimization: If the member is already on the global user list,
			// we don't include them in the member list to avoid redundant lookups.
			if _, ok := reviewerMap[projectMember]; ok {
				continue
			}

			projectMembersList = append(projectMembersList, projectMember)
		}

		allMembersRes, err := conn.PHIDQuery(requests.PHIDQueryRequest{
			PHIDs: projectMembersList,
		})

		if err != nil {
			return nil, err
		}

		// We do the same as for users above, but we assume at this point that all
		// of them are users, not projects.
		for reviewerPHID, queryResult := range allMembersRes {
			reviewerMap[reviewerPHID] = queryResult.Name
		}
	}

	return reviewerMap, nil
}

// lookupSlackMap uses the Phabulous Phabricator extension to lookup the Slack
// account IDs of the reviewers for a revision using their PHIDs. If the
// message is being handled by a non-Slack bot, empty results are returned.
func (c *SummonCommand) lookupSlackMap(
	bot interfaces.Bot,
	conn *gonduit.Conn,
	revision *entities.DifferentialRevision,
) (*responses.PhabulousToSlackResponse, []slack.User, error) {
	// If the extension module is not loaded, we don't attempt to use this
	// functionality.
	if !utilities.HasModule(bot, "extension") {
		return &responses.PhabulousToSlackResponse{}, []slack.User{}, nil
	}

	var slackMap *responses.PhabulousToSlackResponse
	var slackUsers []slack.User
	var err error

	// First, we check that the bot implementation is a Slack bot.
	if sb, ok := bot.(interfaces.SlackBot); ok {
		// Lookup the reviewers using the extension.
		slackMap, err = extensions.PhabulousToSlack(
			conn,
			phabulousRequests.PhabulousToSlackRequest{
				UserPHIDs: revision.Reviewers,
			},
		)

		if err != nil {
			return nil, nil, err
		}

		// Get all Slack users from the Slack API.
		slackUsers, err = sb.GetSlack().GetUsers()

		if err != nil {
			return nil, nil, err
		}
	}

	return slackMap, slackUsers, nil
}

// findOnSlack attempts to find a match for reviewerPHID in slackMap and
// slackUsers. If a match is found, it is returned along true. Otherwise, an
// empty string and false are returned.
func (c *SummonCommand) findOnSlack(
	bot interfaces.Bot,
	slackMap *responses.PhabulousToSlackResponse,
	slackUsers *[]slack.User,
	reviewerPHID string,
) (string, bool) {
	// First, we check that the bot implementation is a Slack bot.
	if _, ok := bot.(interfaces.SlackBot); ok {
		// Next, we check that the reviewer's PHID lookup came back with a result
		// matching it to some Slack ID.
		if slackUserInfo, ok := (*slackMap)[reviewerPHID]; ok {
			var foundUser *slack.User

			// We will go over the list of Slack users and attempt to find a matching
			// account.
			for _, user := range *slackUsers {
				if user.ID == slackUserInfo.AccountID {
					foundUser = &user

					break
				}
			}

			// If we found a user, we return a formatted version of their username to
			// be added to a list of usernames.
			if foundUser != nil {
				formattedName := fmt.Sprintf("@%s :slack:", foundUser.Name)

				return formattedName, true
			}
		}
	}

	return "", false
}
