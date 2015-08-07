package resolvers

import (
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/jacobstr/confer"
)

// CommitResolver resolves phabricator commits to a channel.
type CommitResolver struct {
	Config  *confer.Config            `inject:""`
	Factory *factories.GonduitFactory `inject:""`
}

// Resolve resolves the channel the message about a commit should be posted on.
func (c *CommitResolver) Resolve(name string) (string, error) {
	conduit, err := c.Factory.Make()

	if err != nil {
		return "", err
	}

	commits, err := conduit.DiffusionQueryCommits(
		requests.DiffusionQueryCommitsRequest{
			Names: []string{name},
		},
	)

	if err != nil {
		return "", err
	}

	commit := commits.Data[commits.IdentifierMap[name]]

	reposPtr, err := conduit.RepositoryQuery(requests.RepositoryQueryRequest{
		PHIDs: []string{commit.RepositoryPHID},
		Order: "newest",
	})

	if err != nil {
		return "", err
	}

	repos := *reposPtr
	channelMap := c.Config.GetStringMapString("channels.repositories")

	if channelName, ok := channelMap[repos[0].Callsign]; ok == true {
		return channelName, nil
	}

	return "", nil
}
