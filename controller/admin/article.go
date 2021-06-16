package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"news/bin"
	"news/model"
	"news/resources"
	"os"
)

func RenderHtml(ctx *gin.Context) {
	var articleModels []model.Article
	resources.Db.Find(&articleModels)
	for _, articleModel := range articleModels {
		bin.CacheTemplate(articleModel)
		fmt.Println("success")
	}
}

func Init(ctx *gin.Context) () {
	_ = os.RemoveAll("./cache")
	bin.FetchList()
	ctx.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}

func FetchLatestArticle(ctx *gin.Context) {
	bin.FetchLatestArticle()

	ctx.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}

func FetchArticle(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "Need url",
		})
		return
	}
	bin.FetchFlow(url)
	ctx.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}
