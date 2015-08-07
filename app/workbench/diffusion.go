package workbench

import (
	"github.com/codegangsta/cli"
	"github.com/davecgh/go-spew/spew"
	"github.com/etcinit/gonduit"
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/factories"
)

type DiffusionService struct {
	GonduitFactory *factories.GonduitFactory `inject:""`
}

func (d *DiffusionService) QueryCommitsByName(c *cli.Context) {
	var conn *gonduit.Conn
	var err error
	if conn, err = d.GonduitFactory.Make(); err != nil {
		panic(err)
	}

	res, err := conn.DiffusionQueryCommits(requests.DiffusionQueryCommitsRequest{
		Names: []string{c.Args().First()},
	})

	if err != nil {
		panic(err)
	}

	spew.Dump(res)
}

func (d *DiffusionService) QueryRepositoriesByCallsign(c *cli.Context) {
	var conn *gonduit.Conn
	var err error
	if conn, err = d.GonduitFactory.Make(); err != nil {
		panic(err)
	}

	res, err := conn.RepositoryQuery(requests.RepositoryQueryRequest{
		Callsigns: []string{c.Args().First()},
		Order:     "newest",
	})

	if err != nil {
		panic(err)
	}

	spew.Dump(res)
}
