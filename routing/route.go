package routing

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/routing/controller"
	"github.com/jwmwalrus/felice-franz/routing/controller/api/v1"
)

// Route Configures all the API routes
func Route(r *gin.Engine) *gin.Engine {
	r.Static("/css", "./public/css")
	r.Static("/img", "./public/img")
	// r.Static("/scss", "./public/scss")
	r.Static("/js", "./public/js")
	r.Static("/vendor", "./public/vendor")
	// r.StaticFile("/favicon.ico", "./img/favicon.ico")

	r.HTMLRender = ginview.Default()
	controller.Landing(r)

	api.EndPoints(r)

	return r
}
