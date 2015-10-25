package workbench

import (
	"github.com/codegangsta/cli"
	"github.com/etcinit/gonduit"
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/factories"
)

// ManiphestService provides commands for interacting with the Maniphest
// application.
type ManiphestService struct {
	GonduitFactory *factories.GonduitFactory `inject:""`
}

// QueryByPHIDs queries tasks by their PHIDs
func (d *ManiphestService) QueryByPHIDs(c *cli.Context) {
	spewOrPanic(d.GonduitFactory, func(client *gonduit.Conn) (interface{}, error) {
		return client.ManiphestQuery(requests.ManiphestQueryRequest{
			PHIDs: []string(c.Args()),
		})
	})
}

// QueryByIDs queries tasks by their IDs
func (d *ManiphestService) QueryByIDs(c *cli.Context) {
	spewOrPanic(d.GonduitFactory, func(client *gonduit.Conn) (interface{}, error) {
		return client.ManiphestQuery(requests.ManiphestQueryRequest{
			IDs: []string(c.Args()),
		})
	})
}
