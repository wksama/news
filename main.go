package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"html/template"
	"news/bin"
	"news/controller/admin"
	"news/controller/index"
	"news/model"
	"news/resources"
	"news/utils"
)

func main() {
	c := cron.New()
	_, _ = c.AddFunc("0 14-18 * * *", bin.FetchLatestArticle)
	c.Start()
	resources.Init()
	model.Init()

	if !viper.GetBool("app.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"float64ToString": utils.Float64ToString,
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/", index.List)
	r.GET("/date/:date", index.Item)

	adminGroup := r.Group("/admin")
	adminGroup.GET("/page", admin.LatestPage)
	adminGroup.GET("/cache", admin.CacheArticle)
	adminGroup.GET("/latest", admin.FetchLatestArticle)
	adminGroup.GET("/fetch", admin.FetchArticle)

	_ = r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
