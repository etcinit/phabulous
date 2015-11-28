package factories

import (
	"errors"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/gonduit/core"
	"github.com/jacobstr/confer"
)

var (
	// ErrNoCredentials is used whenever credentials are not set for any
	// authentication method.
	ErrNoCredentials = errors.New("Missing conduit credentials")
)

// GonduitFactory provides functions for building Gounduit clients.
type GonduitFactory struct {
	Config *confer.Config `inject:""`
}

// Make builds an instance of a Gonduit client.
func (g *GonduitFactory) Make() (*gonduit.Conn, error) {
	options := &core.ClientOptions{
		InsecureSkipVerify: g.Config.GetBool("misc.ignore-ca"),
	}

	if g.Config.IsSet("conduit.token") {
		options.APIToken = g.Config.GetString("conduit.token")
	} else if g.Config.IsSet("conduit.cert") && g.Config.IsSet("conduit.user") {
		options.Cert = g.Config.GetString("conduit.cert")
		options.CertUser = g.Config.GetString("conduit.user")
	} else {
		return nil, ErrNoCredentials
	}

	// Attempt to create an instance of a connection.
	conduit, err := gonduit.Dial(
		g.Config.GetString("conduit.api"),
		options,
	)
	if err != nil {
		return nil, err
	}

	// If we are using certificate-based configuration, create a session too.
	if g.Config.IsSet("conduit.user") &&
		g.Config.IsSet("conduit.token") == false {
		err = conduit.Connect()

		if err != nil {
			return nil, err
		}
	}

	return conduit, nil
}
