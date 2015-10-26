package resolvers

import (
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/jacobstr/confer"
)

// DifferentialResolver resolves phabricator revisions and diffs to a channel.
type DifferentialResolver struct {
	Config  *confer.Config            `inject:""`
	Factory *factories.GonduitFactory `inject:""`
}

// Resolve resolves the channel the message about a revision should be posted
// on.
func (c *DifferentialResolver) Resolve(phid string) (string, error) {
	conduit, err := c.Factory.Make()

	if err != nil {
		return "", err
	}

	results, err := conduit.DifferentialQuery(
		requests.DifferentialQueryRequest{
			PHIDs: []string{phid},
		},
	)

	if err != nil {
		return "", err
	}

	revision := (*results)[0]

	repos, err := conduit.RepositoryQuery(requests.RepositoryQueryRequest{
		PHIDs: []string{revision.RepositoryPHID},
		Order: "newest",
	})

	if err != nil {
		return "", err
	}

	channelMap := c.Config.GetStringMapString("channels.repositories")

	if channelName, ok := channelMap[(*repos)[0].Callsign]; ok == true {
		return channelName, nil
	}

	return "", nil
}
