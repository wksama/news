package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"html/template"
	"penti/bin"
	"penti/controller/admin"
	"penti/controller/index"
	"penti/model"
	"penti/resources"
	"penti/utils"
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
	adminGroup.GET("/render", admin.RenderHtml)
	adminGroup.GET("/init", admin.Init)
	adminGroup.GET("/latest", admin.FetchLatestArticle)
	adminGroup.GET("/fetch", admin.FetchArticle)

	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
