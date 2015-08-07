package app

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/jacobstr/confer"
)

// ServeService provides the serve command
type ServeService struct {
	Engine *EngineService `inject:""`
	Config *confer.Config `inject:""`
}

// Run starts up the HTTP server
func (s *ServeService) Run(c *cli.Context) {
	fmt.Println("Starting up the server... (a.k.a. coffee time)")

	engine := s.Engine.New()

	// Figure out which port to use
	port := ":" + strconv.Itoa(s.Config.GetInt("server.port"))

	engine.Run(port)

	fmt.Println("✔︎ Done!")
}
