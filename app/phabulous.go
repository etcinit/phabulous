package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/etcinit/phabulous/app/bot"
	"github.com/etcinit/phabulous/app/workbench"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// Phabulous is the root node of the DI graph
type Phabulous struct {
	Config         *confer.Config                   `inject:""`
	Engine         *EngineService                   `inject:""`
	Serve          *ServeService                    `inject:""`
	Slacker        *bot.SlackService                `inject:""`
	SlackWorkbench *workbench.SlackWorkbenchService `inject:""`
	Logger         *logrus.Logger                   `inject:""`
}

// Boot the upper part of the application.
func (p *Phabulous) Boot() {
	p.Logger.Debugln("Booting upper layer")

	p.Slacker.Slack = slack.New(
		p.Config.GetString("slack.token"),
	)
}
