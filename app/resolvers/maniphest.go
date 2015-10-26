package resolvers

import (
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/jacobstr/confer"
)

// TaskResolver resolves Maniphest tasks to Slack channels.
type TaskResolver struct {
	Config  *confer.Config            `inject:""`
	Factory *factories.GonduitFactory `inject:""`
}

// Resolve resolves the channel the message about a task should be posted on.
func (c *TaskResolver) Resolve(phid string) (string, error) {
	conduit, err := c.Factory.Make()

	if err != nil {
		return "", err
	}

	tasks, err := conduit.ManiphestQuery(
		requests.ManiphestQueryRequest{
			PHIDs: []string{phid},
		},
	)

	if err != nil {
		return "", err
	}

	task := tasks.Get(phid)

	if len(task.ProjectPHIDs) < 1 {
		return "", nil
	}

	projectsPtr, err := conduit.ProjectQuery(requests.ProjectQueryRequest{
		PHIDs: []string{task.ProjectPHIDs[0]},
	})

	if err != nil {
		return "", err
	}

	projects := *projectsPtr
	channelMap := c.Config.GetStringMapString("channels.projects")

	if channelName, ok := channelMap[projects.Data[task.ProjectPHIDs[0]].ID]; ok == true {
		return channelName, nil
	}

	return "", nil
}
