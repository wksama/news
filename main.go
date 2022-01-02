package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"news/bin"
	"news/boot"
	"news/router"
)

func main() {
	boot.Init()

	r := gin.Default()

	flag.Parse()
	runType := flag.Arg(0)
	switch runType {
	case "latest":
		bin.FetchLatestArticle()
	case "page":
		bin.FetchFirstPage()
	case "url":
		url := flag.Arg(1)
		bin.FetchAndStore(url, nil)
	case "cache":
		bin.Cache()
	case "serve":
		fallthrough
	default:
		bin.Cache()
		router.InitRoutes(r)
		port := viper.GetString("app.port")
		_ = r.Run(":" + port)
	}
}
