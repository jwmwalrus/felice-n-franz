package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/bnp/logger"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	"github.com/jwmwalrus/felice-n-franz/internal/kafkian"
)

//go:embed version.json
var version []byte

//go:embed web/html/index.html
var indexHTML string

//go:embed public
var staticContent embed.FS

func main() {
	base.Load(version)
	cfg := base.Conf

	gin.DefaultWriter = logger.NewDefaultWriter()
	r := gin.Default()
	r.Use(logger.ToFile())
	r.Use(gin.Recovery())

	route(r)
	fmt.Printf("\n  Running %v on %v\n", base.AppName, cfg.GetURL())
	r.Run(cfg.GetPort())
}

func route(r *gin.Engine) *gin.Engine {
	sc, err := fs.Sub(staticContent, "public")
	onerror.Panic(err)
	r.StaticFS("/static", http.FS(sc))

	t, err := template.New("index.html").Parse(indexHTML)
	r.SetHTMLTemplate(t)
	r.GET("/", index)
	r.GET("/envs/:envname", getEnv)
	r.GET("/ws", webSocket)

	return r
}

func index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":              base.AppName,
			"version":            base.AppVersion.String(),
			"envs":               base.Conf.GetEnvsList(),
			"replayTypes":        kafkian.GetReplayTypesList(),
			"headerPlaceholder":  `e.g.: [{"key": "Content-Type","value": "application/json"}]`,
			"payloadPlaceholder": `e.g.: {"year": 1984, "enemy": "Eastasia"}`,
		},
	)
}

func getEnv(c *gin.Context) {
	envName := c.Param("envname")
	config := base.Conf.GetEnvConfig(envName)
	c.JSON(http.StatusOK, gin.H{
		"name":         config.Name,
		"headerPrefix": config.HeaderPrefix,
		"groups":       config.Groups,
		"topics":       config.Topics,
		"schemas":      config.Schemas,
	})
}

func webSocket(c *gin.Context) {
	kafkian.HandleWS(c.Writer, c.Request)
}
