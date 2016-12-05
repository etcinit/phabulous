package connectors

import (
	"math/rand"
	"regexp"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules/core"
	"github.com/etcinit/phabulous/app/modules/dev"
	"github.com/etcinit/phabulous/app/modules/extension"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// NewBot creates a new instance of a Bot.
func (c *SlackConnector) setupRTM(slackRTM *slack.RTM, slackInfo *slack.Info) {
	c.slackInfo = slackInfo
	c.slackRTM = slackRTM
	c.imChannelIDs = map[string]bool{}
	c.modules = []interfaces.Module{
		&dev.Module{},
		&core.Module{},
		&extension.Module{},
	}

	// Make it easy to lookup if a channel is an IM channel.
	for _, im := range slackInfo.IMs {
		c.imChannelIDs[im.ID] = true
	}

	// Load message handlers
	c.loadHandlers()
}

// HandlerTuple a tuples of a pattern and a handler.
type HandlerTuple struct {
	Pattern *regexp.Regexp
	Handler interfaces.Handler
}

func (b *SlackConnector) mentionRegex(contents string) *regexp.Regexp {
	username := b.slackInfo.User.ID

	return regexp.MustCompile("^<@" + username + ">:?\\s+" + contents + "$")
}

func (b *SlackConnector) loadHandlers() {
	b.handlers = []HandlerTuple{}
	b.imHandlers = []HandlerTuple{}

	for _, module := range b.modules {
		for _, command := range module.GetCommands() {
			for _, rgx := range command.GetMatchers() {
				b.handlers = append(b.handlers, HandlerTuple{
					Pattern: regexp.MustCompile(rgx),
					Handler: command.GetHandler(),
				})
			}

			for _, rgx := range command.GetIMMatchers() {
				b.imHandlers = append(b.imHandlers, HandlerTuple{
					Pattern: regexp.MustCompile(rgx),
					Handler: command.GetHandler(),
				})
			}

			for _, rgx := range command.GetMentionMatchers() {
				b.handlers = append(b.handlers, HandlerTuple{
					Pattern: b.mentionRegex(rgx),
					Handler: command.GetHandler(),
				})
			}
		}
	}
}

// Excuse comes up with an excuse of why something failed.
func (c *SlackConnector) Excuse(m messages.Message, err error) {
	c.logger.Error(err)

	if c.config.GetBool("server.serious") {
		c.Post(
			m.GetChannel(),
			"An error ocurred and I was unable to fulfill your request.",
			messages.IconDefault,
			true,
		)

		return
	}

	excuses := []string{
		"There is some interference right now and I can't fulfill your request.",
		"Oh noes. I messed up.",
		"Whoops. Something went wrong.",
		"1000s lines of code and I still cant get some things right.",
		"[explodes]",
		"Error: WHY U NO WORK?!",
		"OMG! It failed.",
		"such failure. such request.",
		"Oops I did it again...",
		"A cat is walking over my keywpdfahsgasgdadfk kj h",
	}

	c.Post(
		m.GetChannel(),
		excuses[rand.Intn(len(excuses))],
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
