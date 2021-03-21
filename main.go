package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/api/v1"
	"github.com/jwmwalrus/felice-franz/api/v1/middleware"
	"github.com/jwmwalrus/felice-franz/base"
)

func main() {
	base.Load()
	cfg := base.Conf

	gin.DefaultWriter = middleware.NewDefaultWriter()
	r := gin.Default()
	r.Use(middleware.LoggerToFile())
	r.Use(gin.Recovery())

	api.Route(r)
	fmt.Println("\n  Running", base.AppName, "on", cfg.GetURL())
	r.Run(cfg.GetPort())

	base.Unload()

	fmt.Println("Bye", base.OS, "from", filepath.Base(os.Args[0]))
}
