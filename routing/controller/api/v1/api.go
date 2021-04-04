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

		configEndPoints(v1)
		consumersEndPoints(v1)
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
