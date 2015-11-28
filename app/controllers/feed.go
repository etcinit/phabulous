package controllers

import (
	"github.com/etcinit/gonduit/constants"
	"github.com/etcinit/phabulous/app/bot"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/resolvers"
	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
)

// FeedController handles feed webhook routes
type FeedController struct {
	Config       *confer.Config                  `inject:""`
	Slacker      *bot.SlackService               `inject:""`
	Factory      *factories.GonduitFactory       `inject:""`
	Commits      *resolvers.CommitResolver       `inject:""`
	Tasks        *resolvers.TaskResolver         `inject:""`
	Differential *resolvers.DifferentialResolver `inject:""`
}

// Register registers the route handlers for this controller
func (f *FeedController) Register(r *gin.RouterGroup) {
	front := r.Group("/feed")
	{
		front.POST("/receive", f.postReceive)
	}
}

func (f *FeedController) postReceive(c *gin.Context) {
	conduit, err := f.Factory.Make()

	if err != nil {
		panic(err)
	}

	c.Request.ParseForm()

	res, err := conduit.PHIDQuerySingle(
		string(c.Request.PostForm.Get("storyData[objectPHID]")),
	)

	if err != nil {
		panic(err)
	}

	storyText := c.Request.PostForm.Get("storyText")

	if res.URI != "" {
		storyText += " (<" + res.URI + "|More info>)"
	}

	phidType := constants.PhidType(res.Type)
	icon := messages.PhidTypeToIcon(phidType)

	f.Slacker.FeedPost(storyText)

	switch phidType {
	case constants.PhidTypeCommit:
		if channelName, _ := f.Commits.Resolve(res.Name); channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon)
		}
		break
	case constants.PhidTypeTask:
		if channelName, _ := f.Tasks.Resolve(res.PHID); channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon)
		}
		break
	case constants.PhidTypeDifferentialRevision:
		if channelName, _ := f.Differential.Resolve(res.PHID); channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon)
		}
		break
	}

	c.JSON(200, gin.H{
		"status": "success",
		"messages": []string{
			"OK",
		},
	})
}
