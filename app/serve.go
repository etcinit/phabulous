package app

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/etcinit/phabulous/app/bot"
	"github.com/jacobstr/confer"
)

// ServeService provides the serve command
type ServeService struct {
	Engine  *EngineService    `inject:""`
	Config  *confer.Config    `inject:""`
	Logger  *logrus.Logger    `inject:""`
	Slacker *bot.SlackService `inject:""`
	App     *Phabulous        `inject:""`
}

// Run starts up the HTTP server
func (s *ServeService) Run(c *cli.Context) {
	// Boot the upper layers of the app.
	s.App.Boot(c)

	s.Logger.Infoln("Starting up the server... (a.k.a. coffee time)")

	engine := s.Engine.New()

	go s.Slacker.BootRTM()

	// Try to use the hostname specified on the configuration.
	hostname := ""
	if s.Config.IsSet("server.hostname") {
		hostname = s.Config.GetString("server.hostname")
	}

	// Figure out which port to use.
	hostname = hostname + ":" + strconv.Itoa(s.Config.GetInt("server.port"))

	engine.Run(hostname)

	s.Logger.Infoln("✔︎ Done!")
}
