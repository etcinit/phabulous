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

// LoadModules instructs the connector to load the provided modules.
func (c *IRCConnector) LoadModules(modules []interfaces.Module) {
	c.modules = modules

	if c.client != nil {
		c.loadHandlers()
	}
}

// Excuse comes up with an excuse of why something failed.
func (c *IRCConnector) Excuse(m interfaces.Message, err error) {
	c.logger.Error(err)

	c.Post(
		m.GetChannel(),
		messages.GetExcuse(c.config),
		messages.IconDefault,
		true,
	)
}

// GetHandlers returns the regular handlers loaded in this connector.
func (c *IRCConnector) GetHandlers() []interfaces.HandlerTuple {
	return c.handlers
}

// GetIMHandlers returns the IM handlers loaded in this connector.
func (c *IRCConnector) GetIMHandlers() []interfaces.HandlerTuple {
	return c.imHandlers
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
func (c *IRCConnector) GetModules() []interfaces.Module {
	return c.modules
}

// GetUsername returns the username of the user specified.
//
// Since the IRC connector uses usernames as user IDs, this method is just
// an identity function.
func (c *IRCConnector) GetUsername(userID string) (string, error) {
	return userID, nil
}

// GetUsageHandler returns a handler to be used for when no other handlers are
// matched. This handler usually posts some for of help message.
func (c *IRCConnector) GetUsageHandler() interfaces.Handler {
	return func(b interfaces.Bot, m interfaces.Message, matches []string) {
		c.Post(
			m.GetChannel(),
			"Hi. For usage information, type `help`.",
			messages.IconTasks,
			true,
		)
	}
}

// StartTyping notifies the network that the bot is typing.
//
// The IRC protocol does not provide this feature, so include a dummy method to
// satisfy the interface.
func (c *IRCConnector) StartTyping(channel string) {}

// processMessage processes incoming messages events and calls the appropriate
// handlers.
func (c *IRCConnector) processMessage(conn *client.Conn, line *client.Line) {
	nick := c.client.Me().Nick

	message := NewIRCMessage(line, nick)

	processMessage(c, message)
}

// regexBuilder takes the type of a matcher and its pattern and builds a
// regular expression.
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

// loadHandlers uses the handler compiler to load all the handlers into two
// internal groups.
func (c *IRCConnector) loadHandlers() {
	c.handlers, c.imHandlers = modules.CompileHandlers(
		c.modules,
		c.regexBuilder,
	)
}
