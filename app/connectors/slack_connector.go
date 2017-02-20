package connectors

import (
	"github.com/Sirupsen/logrus"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// NewSlackConnector constructs an instance of a SlackConnector.
func NewSlackConnector(
	config *confer.Config,
	gonduitFactory *factories.GonduitFactory,
	logger *logrus.Logger,
) *SlackConnector {
	slack := slack.New(config.GetString("slack.token"))

	connector := &SlackConnector{
		logger:         logger,
		config:         config,
		slack:          slack,
		gonduitFactory: gonduitFactory,
		modules:        []interfaces.Module{},
	}

	return connector
}

// SlackConnector provides a connector service to Slack networks.
type SlackConnector struct {
	logger *logrus.Logger
	config *confer.Config

	slack *slack.Client

	gonduitFactory *factories.GonduitFactory

	slackInfo    *slack.Info
	slackRTM     *slack.RTM
	imChannelIDs map[string]bool
	handlers     []interfaces.HandlerTuple
	imHandlers   []interfaces.HandlerTuple

	modules []interfaces.Module
}

// Boot initializes the connector.
func (c *SlackConnector) Boot() error {
	c.logger.Infoln("Starting RTM handler...")

	rtm := c.slack.NewRTM()
	go rtm.ManageConnection()
	go c.rtmLoop(rtm)

	return nil
}

func (c *SlackConnector) rtmLoop(rtm *slack.RTM) {
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			c.logger.Debugln("RTM Event received.")

			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				c.logger.Debugln("Infos:", ev.Info)
				c.logger.Debugln("Connection counter:", ev.ConnectionCount)

				c.setupRTM(rtm, ev.Info)

				c.logger.Infoln("Bot Slack ID: ", ev.Info.User.ID)

				c.PostOnFeed("Hi! Phabulous v3 reporting for duty!")
			case *slack.MessageEvent:
				c.logger.Debugf("Message: %v\n", ev)
				c.processMessage(ev)

			case *slack.IMOpenEvent:
				c.processIMOpen(ev)

			case *slack.PresenceChangeEvent:
				c.logger.Debugf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				c.logger.Debugf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				c.logger.Errorf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				c.logger.Errorln("Invalid credentials.")

				break Loop

			case *slack.ConnectingEvent:
				c.logger.Infof(
					"Attempting to connect to Slack (Attempt #%d)",
					ev.Attempt,
				)
			case *slack.ConnectionErrorEvent:
				c.logger.Error(
					"Unable to connect/authenticate with Slack. ",
					"Check the bot's credentials.",
				)

			default:
				// Ignore other events..
			}
		}
	}

	c.logger.Warnln("RTM handler has stopped.")
}
