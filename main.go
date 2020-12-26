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
	"penti/utils"
)

func main() {
	model.InitDb()

	c := cron.New()
	_, _ = c.AddFunc("0 * * * *", bin.FetchLatestArticle)
	c.Start()

	if !viper.GetBool("app.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	//r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"float64ToString": utils.Float64ToString,
	})

	r.LoadHTMLGlob("templates/*" )

	r.GET("/", index.List)
	r.GET("/date/:date", index.Item)

	r.GET("/render", admin.RenderHtml)
	r.GET("/init", admin.Init)

	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
