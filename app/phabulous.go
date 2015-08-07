package app

import (
	"github.com/etcinit/phabulous/app/slacker"
	"github.com/etcinit/phabulous/app/workbench"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// Phabulous is the root node of the DI graph
type Phabulous struct {
	Config    *confer.Config              `inject:""`
	Engine    *EngineService              `inject:""`
	Serve     *ServeService               `inject:""`
	Slacker   *slacker.SlackService       `inject:""`
	Diffusion *workbench.DiffusionService `inject:""`
}

// Boot the upper part of the application.
func (p *Phabulous) Boot() {
	p.Slacker.Slack = slack.New(
		p.Config.GetString("slack.token"),
	)
}
