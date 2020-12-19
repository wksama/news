package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"html/template"
	"penti/bin"
	"penti/model"
	"penti/utils"
)

func main() {
	model.InitDb()

	c := cron.New()
	_, _ = c.AddFunc("* * * * *", bin.FetchLatestArticle)
	c.Start()

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"float64ToString": utils.Float64ToString,
	})

	r.LoadHTMLGlob("templates/*" )

	r.GET("/", List)
	r.GET("/render", bin.RenderHtml)
	r.GET("/init", Init)
	r.GET("/today", Today)
	r.GET("/date/:date", Item)
	//r.GET("/:date", List)

	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
