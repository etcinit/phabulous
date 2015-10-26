package controllers

import (
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/resolvers"
	"github.com/etcinit/phabulous/app/slacker"
	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// FeedController handles feed webhook routes
type FeedController struct {
	Config       *confer.Config                  `inject:""`
	Slacker      *slacker.SlackService           `inject:""`
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

	res, err := conduit.PHIDQuerySingle(string(c.Request.PostForm.Get("storyData[objectPHID]")))

	if err != nil {
		panic(err)
	}

	storyText := c.Request.PostForm.Get("storyText")

	if res.URI != "" {
		storyText += " (<" + res.URI + "|More info>)"
	}

	f.Slacker.Slack.PostMessage(
		f.Config.GetString("channels.feed"),
		storyText,
		slack.PostMessageParameters{
			Username: f.Config.GetString("slack.username"),
			IconURL:  "http://i.imgur.com/7Hzgo9Y.png",
		},
	)

	switch res.Type {
	case "CMIT":
		if channelName, _ := f.Commits.Resolve(res.Name); channelName != "" {
			f.Slacker.Slack.PostMessage(
				channelName,
				storyText,
				slack.PostMessageParameters{
					Username: f.Config.GetString("slack.username"),
					IconURL:  "http://i.imgur.com/v8ReRKx.png",
				},
			)
		}
		break
	case "TASK":
		if channelName, _ := f.Tasks.Resolve(res.PHID); channelName != "" {
			f.Slacker.Slack.PostMessage(
				channelName,
				storyText,
				slack.PostMessageParameters{
					Username: f.Config.GetString("slack.username"),
					IconURL:  "http://i.imgur.com/jD7rf9x.png",
				},
			)
		}
		break
	case "DREV":
		if channelName, _ := f.Differential.Resolve(res.PHID); channelName != "" {
			f.Slacker.Slack.PostMessage(
				channelName,
				storyText,
				slack.PostMessageParameters{
					Username: f.Config.GetString("slack.username"),
					IconURL:  "http://i.imgur.com/NiPouYj.png",
				},
			)
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
