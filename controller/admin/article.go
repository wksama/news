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

func CacheArticle(ctx *gin.Context) {
	var page = 1
	var size = 20
	for true {
		var offset = (page - 1) * size
		var articleModels []model.Article
		resources.Db.Offset(offset).Limit(size).Order("date DESC").Find(&articleModels)
		for _, articleModel := range articleModels {
			bin.CacheFlow(articleModel)
			fmt.Println(articleModel.FullTitle + ": cache success")
		}
		if len(articleModels) < size {
			break
		}
		page++
	}
}

func LatestPage(ctx *gin.Context) {
	_ = os.RemoveAll("./cache")
	bin.FetchFirstPage()
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func FetchLatestArticle(ctx *gin.Context) {
	bin.FetchLatestArticle()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
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
