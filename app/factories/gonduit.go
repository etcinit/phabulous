package factories

import (
	"github.com/etcinit/gonduit"
	"github.com/jacobstr/confer"
)

// GonduitFactory provides functions for building Gounduit clients.
type GonduitFactory struct {
	Config *confer.Config `inject:""`
}

// Make builds an instance of a Gonduit client.
func (g *GonduitFactory) Make() (*gonduit.Conn, error) {
	conduit, err := gonduit.Dial(
		g.Config.GetString("conduit.api"),
		&gonduit.ClientOptions{
			InsecureSkipVerify: g.Config.GetBool("misc.ignore-ca"),
		},
	)

	if err != nil {
		return nil, err
	}

	err = conduit.Connect(
		g.Config.GetString("conduit.user"),
		g.Config.GetString("conduit.cert"),
	)

	if err != nil {
		return nil, err
	}

	return conduit, nil
}
