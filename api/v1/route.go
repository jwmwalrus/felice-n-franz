package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/base"
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

func index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"title": base.AppName,
		},
	)
}

func about(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/about.html",
		gin.H{},
	)
}

func contact(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"views/contact.html",
		gin.H{
			"title": "Contact",
		},
	)
}
