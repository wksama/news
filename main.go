package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"
	"news/bin"
	"news/controller/admin"
	"news/controller/index"
	"news/model"
	"news/resources"
	"news/utils"
	"os"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe(":1234", nil)
	}()

	// manually set time zone
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

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

	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.gohtml", nil)
	})

	_ = r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
