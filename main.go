package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	"github.com/jwmwalrus/felice-n-franz/internal/kafkian"
)

//go:embed web/html/index.html
var indexHTML string

//go:embed public
var staticContent embed.FS

func main() {
	base.Load()
	cfg := base.Conf

	gin.DefaultWriter = bnp.NewDefaultWriter()
	r := gin.Default()
	r.Use(bnp.LoggerToFile())
	r.Use(gin.Recovery())

	route(r)
	fmt.Println("\n  Running", base.AppName, "on", cfg.GetURL())
	r.Run(cfg.GetPort())

	base.Unload()

	fmt.Println("Bye", base.OS, "from", filepath.Base(os.Args[0]))
}

func route(r *gin.Engine) *gin.Engine {
	sc, err := fs.Sub(staticContent, "public")
	bnp.PanicOnError(err)
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
			"envs":               base.Conf.GetEnvsList(),
			"headerPlaceholder":  `e.g.: [{"key": "Content-Type","value": "application/json"}]`,
			"payloadPlaceholder": `e.g.: {"year": 1984, "enemy": "Eastasia"}`,
		},
	)
}

func getEnv(c *gin.Context) {
	envName := c.Param("envname")
	config := base.Conf.GetEnvConfig(envName)
	c.JSON(http.StatusOK, gin.H{
		"name":          config.Name,
		"headerPrefix":  config.HeaderPrefix,
		"maxTailOffset": base.Conf.MaxTailOffset,
		"groups":        config.Groups,
		"topics":        config.Topics,
		"schemas":       config.Schemas,
	})
}

func webSocket(c *gin.Context) {
	kafkian.HandleWS(c.Writer, c.Request)
}
