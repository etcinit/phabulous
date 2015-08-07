package controllers

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/slacker"
	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
)

// FeedController handles feed webhook routes
type FeedController struct {
	Config  *confer.Config        `inject:""`
	Slacker *slacker.SlackService `inject:""`
}

// Register registers the route handlers for this controller
func (f *FeedController) Register(r *gin.RouterGroup) {
	front := r.Group("/feed")
	{
		front.POST("/receive", f.postReceive)
	}
}

func (f *FeedController) postReceive(c *gin.Context) {
	conduit, err := gonduit.Dial(f.Config.GetString("conduit.api"))

	if err != nil {
		panic(err)
	}

	err = conduit.Connect(
		f.Config.GetString("conduit.user"),
		f.Config.GetString("conduit.cert"),
	)

	if err != nil {
		panic(err)
	}

	c.Request.ParseForm()
	spew.Dump(c.Request.PostForm)

	//res, err := conduit.PHIDQuerySingle(c.Request.Get(""))

	c.JSON(200, gin.H{
		"status": "success",
		"messages": []string{
			"OK",
		},
	})
}
