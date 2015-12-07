package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/etcinit/phabulous/app/bot"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// Phabulous is the root node of the DI graph
type Phabulous struct {
	Config  *confer.Config    `inject:""`
	Engine  *EngineService    `inject:""`
	Serve   *ServeService     `inject:""`
	Slacker *bot.SlackService `inject:""`
	Logger  *logrus.Logger    `inject:""`
}

// Boot the upper part of the application.
func (p *Phabulous) Boot(c *cli.Context) {
	if c.GlobalString("config") != "" {
		err := p.Config.ReadPaths(c.GlobalString("config"))

		if err != nil {
			p.Logger.Panic(err)
		}

		p.Logger.Infoln(
			"Loaded alternate configuration file from: " +
				c.GlobalString("config") + ".",
		)
	}

	if p.Config.GetBool("server.debug") {
		p.Logger.Level = logrus.DebugLevel
		p.Logger.Debugln("Logger is debug level.")
	} else {
		p.Logger.Level = logrus.WarnLevel
	}

	p.Slacker.Slack = slack.New(
		p.Config.GetString("slack.token"),
	)

	p.Logger.Debugln("Booted upper layer.")
}
