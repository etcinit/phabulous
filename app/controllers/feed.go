package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/etcinit/gonduit/constants"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/resolvers"
	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
)

// FeedController handles feed webhook routes
type FeedController struct {
	Config       *confer.Config                  `inject:""`
	Factory      *factories.GonduitFactory       `inject:""`
	Commits      *resolvers.CommitResolver       `inject:""`
	Tasks        *resolvers.TaskResolver         `inject:""`
	Differential *resolvers.DifferentialResolver `inject:""`
	Logger       *logrus.Logger                  `inject:""`

	Connector interfaces.Connector
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
		f.Logger.Error("ERROR making factory: ", err)
		return
	}

	c.Request.ParseForm()

	f.Logger.Debug(c.Request.PostForm.Encode())

	inString := string(c.Request.PostForm.Get("storyData[objectPHID]"))
	res, err := conduit.PHIDQuerySingle(inString)

	if err != nil {
		f.Logger.Error("ERROR with input ", inString, ": ", err)
		return
	}

	storyText := c.Request.PostForm.Get("storyText")

	if res.URI != "" {
		storyText += " (<" + res.URI + "|More info>)"
	}

	phidType := constants.PhidType(res.Type)
	icon := messages.PhidTypeToIcon(phidType)

	if typeAllowed(f.Config, res.Type) {
		f.Connector.PostOnFeed(storyText)
	}

	switch phidType {
	case constants.PhidTypeCommit:
		channelName, err := f.Commits.Resolve(res.Name)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Connector.Post(channelName, storyText, icon, false)
		}

		// Support "all" channel.
		channelMap := f.Config.GetStringMapString("channels.repositories")

		if channelName, ok := channelMap["all"]; ok == true {
			f.Connector.Post(channelName, storyText, icon, false)
		}
		break
	case constants.PhidTypeTask:
		channelNames, err := f.Tasks.Resolve(res.PHID)
		if err != nil {
			f.Logger.Error(err)
		}
	    for _,channelName := range channelNames {
			f.Connector.Post(channelName, storyText, icon, false)
		}
		break
	case constants.PhidTypeDifferentialRevision:
		channelName, err := f.Differential.Resolve(res.PHID)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Connector.Post(channelName, storyText, icon, false)
		}

		// Support "all" channel.
		channelMap := f.Config.GetStringMapString("channels.repositories")

		if channelName, ok := channelMap["all"]; ok == true {
			f.Connector.Post(channelName, storyText, icon, false)
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

// typeAllowed returns whether a feed event should be processed in the main
// feed channel.
func typeAllowed(config *confer.Config, currentType string) bool {
	allowedTypes := config.GetStringSlice("channels.feedTypes")

	// If no filters are specified, everything is allowed.
	if len(allowedTypes) == 0 {
		return true
	}

	for _, eventType := range allowedTypes {
		if eventType == currentType {
			return true
		}
	}

	return false
}
