package app

import (
	"github.com/etcinit/phabulous/app/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
)

// EngineService provides the API engine
type EngineService struct {
	Front  controllers.FrontController `inject:"inline"`
	Feed   controllers.FeedController  `inject:"inline"`
	Config *confer.Config              `inject:""`
}

// New creates a new instance of an API engine
func (e *EngineService) New() *gin.Engine {
	if e.Config.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	e.Front.Register(router)

	v1 := router.Group("/v1")
	{
		e.Feed.Register(v1)
	}

	return router
}
