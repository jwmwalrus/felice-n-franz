package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// EndPoints returns the defined endpoints
func EndPoints(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", alive)
		v1.POST("/off", off)

		v1.GET("/config", getGlobalConfig)

		v1.GET("/env", getEnvsList)
		v1.GET("/env/:envname/config", getEnvConfig)
		v1.GET("/env/:envname/topic", getEnvTopics)
		v1.GET("/env/:envname/group", getEnvGroups)
		v1.GET("/env/:envname/group/:groupname/topics", getEnvGroupTopics)
	}
}

func alive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
}

func off(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}
