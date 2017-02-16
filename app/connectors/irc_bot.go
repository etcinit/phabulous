package connectors

import (
	"regexp"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/fluffle/goirc/client"
	"github.com/jacobstr/confer"
)

func (c *IRCConnector) LoadModules(modules []interfaces.Module) {
	c.modules = modules

	if c.client != nil {
		c.loadHandlers()
	}
}

// ProcessMessage processes incoming messages events and calls the appropriate
// handlers.
func (c *IRCConnector) processMessage(conn *client.Conn, line *client.Line) {
	nick := c.client.Me().Nick

	// Ignore messages from the bot itself.
	if line.Nick == nick {
		return
	}

	// If the message is an IM, use IM handlers.
	if line.Target() == line.Nick {
		handled := false

		for _, tuple := range c.imHandlers {
			if result := tuple.Pattern.FindStringSubmatch(line.Text()); result != nil {
				go tuple.Handler(c, messages.NewIRCMessage(line), result)

				handled = true
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go c.HandleUsage(messages.NewIRCMessage(line), []string{})
		}

		return
	}

	for _, tuple := range c.handlers {
		if result := tuple.Pattern.FindStringSubmatch(line.Text()); result != nil {
			go tuple.Handler(c, messages.NewIRCMessage(line), result)
		}
	}
}

func (c *IRCConnector) regexBuilder(
	matcherType modules.MatcherType,
	pattern string,
) *regexp.Regexp {
	if matcherType != modules.MentionMatcherType {
		return regexp.MustCompile(pattern)
	}

	username := c.client.Me().Nick

	return regexp.MustCompile("^@" + username + ":? " + pattern + "$")
}

func (c *IRCConnector) loadHandlers() {
	c.handlers, c.imHandlers = modules.CompileHandlers(
		c.modules,
		c.regexBuilder,
	)
}

// Excuse comes up with an excuse of why something failed.
func (c *IRCConnector) Excuse(m messages.Message, err error) {
	c.logger.Error(err)

	c.Post(
		m.GetChannel(),
		messages.GetExcuse(c.config),
		messages.IconDefault,
		true,
	)
}

// GetGonduit gets an instance of a gonduit client.
func (c *IRCConnector) GetGonduit() (*gonduit.Conn, error) {
	return c.gonduitFactory.Make()
}

// GetConfig returns an instance of the configuration store.
func (c *IRCConnector) GetConfig() *confer.Config {
	return c.config
}

// GetModules returns the modules used in this bot.
func (b *IRCConnector) GetModules() []interfaces.Module {
	return b.modules
}

// HandleUsage shows usage tip.
func (c *IRCConnector) HandleUsage(m messages.Message, matches []string) {
	c.Post(
		m.GetChannel(),
		"Hi. For usage information, type `help`.",
		messages.IconTasks,
		true,
	)
}

func (c *IRCConnector) GetUsername(userId string) (string, error) {
	return userId, nil
}

// StartTyping notify Slack that the bot is "typing".
func (b *IRCConnector) StartTyping(channel string) {}
