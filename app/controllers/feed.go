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
	Config  *confer.Config            `inject:""`
	Slacker *slacker.SlackService     `inject:""`
	Factory *factories.GonduitFactory `inject:""`
	Commits *resolvers.CommitResolver `inject:""`
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

	f.Slacker.Slack.PostMessage(
		f.Config.GetString("channels.feed"),
		c.Request.PostForm.Get("storyText"),
		slack.PostMessageParameters{
			Username: f.Config.GetString("slack.username"),
		},
	)

	if res.Type == "CMIT" {
		if channelName, _ := f.Commits.Resolve(res.Name); channelName != "" {
			f.Slacker.Slack.PostMessage(
				channelName,
				c.Request.PostForm.Get("storyText"),
				slack.PostMessageParameters{
					Username: f.Config.GetString("slack.username"),
				},
			)
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"messages": []string{
			"OK",
		},
	})
}
