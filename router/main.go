package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"news/controller/index"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/", index.List)
	r.GET("/date/:date", index.Item)

	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.gohtml", nil)
	})
}
