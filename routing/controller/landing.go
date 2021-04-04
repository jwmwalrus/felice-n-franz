package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/base"
)

// Landing defines homepage-related controllers
func Landing(r *gin.Engine) {
	r.GET("/", index)
	r.GET("/about", about)
}

func index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title": base.AppName,
			"envs":  base.Conf.GetEnvsList(),
		},
	)
}

func about(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"about.html",
		gin.H{},
	)
}
