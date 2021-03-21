package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Route Configures all the API routes
func Route(r *gin.Engine) *gin.Engine {
	r.Static("/css", "./public/css")
	r.Static("/img", "./public/img")
	// r.Static("/scss", "./public/scss")
	r.Static("/vendor", "./public/vendor")
	r.Static("/js", "./public/js")
	// r.StaticFile("/favicon.ico", "./img/favicon.ico")

	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"alive": true})
		})
		v1.POST("/off", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"alive": true})
			go func() {
				time.Sleep(5 * time.Second)
				os.Exit(0)
			}()
		})
	}

	return r
}
