package resolvers

import (
	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/jacobstr/confer"
	"github.com/Sirupsen/logrus"
)

// TaskResolver resolves Maniphest tasks to Slack channels.
type TaskResolver struct {
	Config  *confer.Config            `inject:""`
	Factory *factories.GonduitFactory `inject:""`
    Logger  *logrus.Logger            `inject:""`
}

// Resolve resolves the channel the message about a task should be posted on.
func (c *TaskResolver) Resolve(phid string) ([]string, error) {
	conduit, err := c.Factory.Make()
    var retVal []string
	if err != nil {
		return retVal, err
	}

	tasks, err := conduit.ManiphestQuery(
		requests.ManiphestQueryRequest{
			PHIDs: []string{phid},
		},
	)

	if err != nil {
		return retVal, err
	}

	task := tasks.Get(phid)

	if len(task.ProjectPHIDs) < 1 {
		return retVal, nil
	}
	channelMap := c.Config.GetStringMapString("channels.projects")
    matches := 0
    for _, projectPHID := range task.ProjectPHIDs {
        projects, err := conduit.ProjectQuery(requests.ProjectQueryRequest{
		    PHIDs: []string{projectPHID},
	    })
    	if err != nil {
	    	return retVal, err
    	}
        if channelName, ok := channelMap[projects.Data[projectPHID].ID]; ok == true {
		    retVal = append(retVal, channelName)
                    matches += 1
	    }
    }
    if (matches == 0) {
        retVal = append(retVal, "")
    }
    return retVal, nil
}
