package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"penti/bin"
	"penti/model"
	"penti/utils"
)

func RenderHtml(ctx *gin.Context) {
	var articleModels []model.Article
	model.Db.Find(&articleModels)
	for _,articleModel := range articleModels {
		articleStruct := utils.Model2Article(articleModel)
		htmlBuffer := utils.RenderHtml(articleStruct)
		file := fmt.Sprintf("%s/%s", utils.GetSaveDir(articleModel), utils.GetSaveName(articleModel))
		_ = ioutil.WriteFile(file, htmlBuffer.Bytes(), 0777)
		fmt.Println("success")
	}
}

func Init(ctx *gin.Context)()  {
	_ = os.RemoveAll("./cache")
	bin.FetchList()
	ctx.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}
