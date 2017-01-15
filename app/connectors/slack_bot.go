package connectors

import (
	"regexp"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

func (c *SlackConnector) setupRTM(slackRTM *slack.RTM, slackInfo *slack.Info) {
	c.slackInfo = slackInfo
	c.slackRTM = slackRTM
	c.imChannelIDs = map[string]bool{}

	// Make it easy to lookup if a channel is an IM channel.
	for _, im := range slackInfo.IMs {
		c.imChannelIDs[im.ID] = true
	}

	// Reload handlers
	c.loadHandlers()
}

func (c *SlackConnector) LoadModules(modules []interfaces.Module) {
	c.modules = modules

	// Reload handlers
	if c.slackRTM != nil {
		c.loadHandlers()
	}
}

func (b *SlackConnector) regexBuilder(
	matcherType modules.MatcherType,
	pattern string,
) *regexp.Regexp {
	if matcherType != modules.MentionMatcherType {
		return regexp.MustCompile(pattern)
	}

	username := b.slackInfo.User.ID

	return regexp.MustCompile("^<@" + username + ">:? " + pattern + "$")
}

func (b *SlackConnector) loadHandlers() {
	b.handlers, b.imHandlers = modules.CompileHandlers(
		b.modules,
		b.regexBuilder,
	)
}

// Excuse comes up with an excuse of why something failed.
func (c *SlackConnector) Excuse(m messages.Message, err error) {
	c.logger.Error(err)

	c.Post(
		m.GetChannel(),
		messages.GetExcuse(c.config),
		messages.IconDefault,
		true,
	)
}

// ProcessIMOpen handles IM open events.
func (c *SlackConnector) processIMOpen(ev *slack.IMOpenEvent) {
	c.imChannelIDs[ev.Channel] = true
}

// ProcessMessage processes incoming messages events and calls the appropriate
// handlers.
func (c *SlackConnector) processMessage(ev *slack.MessageEvent) {
	// Ignore messages from the bot itself.
	if ev.User == c.slackInfo.User.ID {
		return
	}

	// If the message is an IM, use IM handlers.
	if _, ok := c.imChannelIDs[ev.Channel]; ok {
		handled := false

		for _, tuple := range c.imHandlers {
			if result := tuple.Pattern.FindStringSubmatch(ev.Text); result != nil {
				go tuple.Handler(c, messages.NewSlackMessage(ev), result)

				handled = true
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go c.HandleUsage(messages.NewSlackMessage(ev), []string{})
		}

		return
	}

	for _, tuple := range c.handlers {
		if result := tuple.Pattern.FindStringSubmatch(ev.Text); result != nil {
			go tuple.Handler(c, messages.NewSlackMessage(ev), result)
		}
	}
}

// GetModules returns the modules used in this bot.
func (b *SlackConnector) GetModules() []interfaces.Module {
	return b.modules
}

// StartTyping notify Slack that the bot is "typing".
func (b *SlackConnector) StartTyping(channel string) {
	b.slackRTM.SendMessage(b.slackRTM.NewTypingMessage(channel))
}

// GetGonduit gets an instance of a gonduit client.
func (c *SlackConnector) GetGonduit() (*gonduit.Conn, error) {
	return c.gonduitFactory.Make()
}

// GetConfig returns an instance of the configuration store.
func (c *SlackConnector) GetConfig() *confer.Config {
	return c.config
}

func (c *SlackConnector) GetSlack() *slack.Client {
	return c.slack
}

// HandleUsage shows usage tip.
func (c *SlackConnector) HandleUsage(m messages.Message, matches []string) {
	c.Post(
		m.GetChannel(),
		"Hi. For usage information, type `help`.",
		messages.IconTasks,
		true,
	)
}

func (c *SlackConnector) GetUsername(userId string) (string, error) {
	userInfo, err := c.slack.GetUserInfo(userId)
	if err != nil {
		return "", err
	}

	return userInfo.Name, nil
}
