package routing

import (
	"net/http"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/base"
)

// Route Configures all the API routes
func Route(r *gin.Engine) *gin.Engine {
	r.Static("/static", "./public/static")
	// r.StaticFile("/favicon.ico", "./img/favicon.ico")

	gv := goview.DefaultConfig
	gv.Root = "assets/html"
	r.HTMLRender = ginview.New(gv)
	r.GET("/", index)
	r.GET("/envs/:envname", getEnv)
	r.GET("/ws", socket)

	return r
}

func index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": base.AppName,
			"envs":  base.Conf.GetEnvsList(),
		},
	)
}

func getEnv(c *gin.Context) {
	envName := c.Param("envname")
	config := base.Conf.GetEnvConfig(envName)
	c.JSON(http.StatusOK, gin.H{
		"name":   config.Name,
		"groups": config.Groups,
		"topics": config.Topics,
	})
}

func socket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
