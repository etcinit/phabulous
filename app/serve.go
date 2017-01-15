package app

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/jacobstr/confer"
)

// ServeService provides the serve command
type ServeService struct {
	Engine *EngineService `inject:""`
	Config *confer.Config `inject:""`
	Logger *logrus.Logger `inject:""`
	App    *Phabulous     `inject:""`
}

// Run boots the application and starts up the API server.
func (s *ServeService) Run(c *cli.Context) {
	// Boot the upper layers of the app.
	s.App.Boot(c)

	// We pass down the connector manager as the feed's connector. This will
	// allow it to post to all services the bot is configured to use.
	s.Engine.Feed.Connector = s.App.ConnectorManager

	// Finally, we boot the connector manager.
	s.App.ConnectorManager.Boot()

	// Start the API server so we can receive feed events.
	s.runServer()

	s.Logger.Infoln("✔︎ Done!")
}

// runServer sets up the HTTP server and begins listening.
func (s *ServeService) runServer() {
	s.Logger.Infoln("Starting up the API server...")

	engine := s.Engine.New()

	// Try to use the hostname specified on the configuration.
	hostname := ""
	if s.Config.IsSet("server.hostname") {
		hostname = s.Config.GetString("server.hostname")
	}

	// Figure out which port to use.
	hostname = hostname + ":" + strconv.Itoa(s.Config.GetInt("server.port"))

	s.Logger.Debugf("API Server Hostname: %s", hostname)

	engine.Run(hostname)
}
