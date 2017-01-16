package connectors

import (
	"crypto/tls"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/interfaces"
	irc "github.com/fluffle/goirc/client"
	"github.com/jacobstr/confer"
)

type IRCConnector struct {
	config         *confer.Config
	gonduitFactory *factories.GonduitFactory
	logger         *logrus.Logger

	handlers   []interfaces.HandlerTuple
	imHandlers []interfaces.HandlerTuple
	modules    []interfaces.Module

	client *irc.Conn
}

// NewIRCConnector constructs a new instance of an IRCConnector.
func NewIRCConnector(
	config *confer.Config,
	gonduitFactory *factories.GonduitFactory,
	logger *logrus.Logger,
) *IRCConnector {
	connector := IRCConnector{}

	connector.config = config
	connector.gonduitFactory = gonduitFactory
	connector.logger = logger

	return &connector
}

func (c *IRCConnector) Boot() error {
	c.logger.Info("Booting IRC connector...")

	hostname := c.config.GetString("irc.hostname")
	port := c.config.GetInt("irc.port")

	cfg := irc.NewConfig(c.config.GetString("irc.nick"))

	cfg.SSL = c.config.GetBool("irc.tls")
	cfg.SSLConfig = &tls.Config{ServerName: hostname}
	cfg.Server = fmt.Sprintf("%s:%d", hostname, port)
	cfg.NewNick = func(n string) string { return n + "^" }

	c.client = irc.Client(cfg)

	c.client.HandleFunc(
		irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			c.logger.Infof("Connected to IRC server: %s:%d", hostname, port)

			c.PostOnFeed("Hi! Phabulous v2 reporting for duty!")

			c.joinConfiguredChannels()

			c.loadHandlers()
		},
	)

	c.client.HandleFunc(
		irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			c.logger.Warnf(
				"Disconnected from IRC server: %s:%d",
				hostname,
				port,
			)
		},
	)

	c.client.HandleFunc(
		irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			spew.Dump(line.Text())
			c.processMessage(conn, line)
		},
	)

	c.client.EnableStateTracking()

	if err := c.client.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
		return err
	}

	return nil
}

func (c *IRCConnector) joinConfiguredChannels() {
	feedChannel := c.config.GetString("channels.feed")
	repositoryChannels := c.config.GetStringMapString("channels.repositories")
	projectChannels := c.config.GetStringMapString("channels.projects")

	if feedChannel != "" {
		c.client.Join(feedChannel)
	}

	channels := map[string]bool{}

	for _, channel := range repositoryChannels {
		channels[channel] = true
	}

	for _, channel := range projectChannels {
		channels[channel] = true
	}

	for channel, _ := range channels {
		c.client.Join(channel)
	}
}
