package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/etcinit/phabulous/app/connectors"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/jacobstr/confer"
)

// Phabulous is the root node of the DI graph
type Phabulous struct {
	Config           *confer.Config            `inject:""`
	Engine           *EngineService            `inject:""`
	Serve            *ServeService             `inject:""`
	Logger           *logrus.Logger            `inject:""`
	GonduitFactory   *factories.GonduitFactory `inject:""`
	ConnectorManager *ConnectorManager         `inject:""`
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

	if p.Config.GetBool("slack.enable") {
		p.ConnectorManager.RegisterConnector(connectors.NewSlackConnector(
			p.Config,
			p.GonduitFactory,
			p.Logger,
		))
	}

	if p.Config.GetBool("irc.enable") {
		p.ConnectorManager.RegisterConnector(connectors.NewIRCConnector(
			p.Config,
			p.GonduitFactory,
			p.Logger,
		))
	}

	p.ConnectorManager.LoadModules(
		modules.NewModuleFactory(p.Config, p.Logger).Make(),
	)

	p.Logger.Debugln("Booted upper layer.")
}
