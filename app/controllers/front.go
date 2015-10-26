package controllers

import "github.com/gin-gonic/gin"

// FrontController handles main routes
type FrontController struct{}

// Register registers the route handlers for this controller
func (f *FrontController) Register(r *gin.Engine) {
	front := r.Group("/")
	{
		front.GET("/", f.getIndex)
		front.GET("/healthcheck", f.getHealthCheck)
	}

	r.NoRoute(f.getNotFound)
}

func (f *FrontController) getIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "success",
		"messages": []string{
			"Welcome to the Phabulous API",
		},
		"version": "2.1.0-alpha1",
	})
}

func (f *FrontController) getHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":   "success",
		"messages": []string{"OK"},
	})
}

func (f *FrontController) getNotFound(c *gin.Context) {
	c.JSON(404, gin.H{
		"status":   "error",
		"messages": []string{"Path not found"},
	})
}
