package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/datatypes"
	"io/ioutil"
	"net/http"
	"os"
	"penti/bin"
	"penti/model"
	"penti/utils"
	"time"
)

const PAGE_SIZE = 20000000

func List(ctx *gin.Context) {
	res := model.Rdb.ZRevRangeWithScores(ctx, "articleList",0, PAGE_SIZE)
	if len(res.Val()) == 0 {
		var modelList []model.Article
		model.Db.Select("full_title, date").Order("date desc").Find(&modelList)

		for _,item := range modelList {
			model.Rdb.ZAdd(ctx, "articleList", &redis.Z{
				Score: utils.DateToFloat64(item.Date),
				Member: item.FullTitle,
			})
		}
		res = model.Rdb.ZRevRangeWithScores(ctx, "articleList",0, PAGE_SIZE)
	}
	list := res.Val()
	ctx.HTML(http.StatusOK, "list.gohtml", list)
}

func Item(ctx *gin.Context) {
	dateStr := ctx.Param("date")

	dateTime, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	var articleModel model.Article
	model.Db.Where("date = ?", datatypes.Date(dateTime)).First(&articleModel)

	outputPage(articleModel, ctx)
	return
}

func Today(ctx *gin.Context) {
	var articleModel model.Article
	model.Db.Where("date = ?", datatypes.Date(time.Now())).First(&articleModel)

	articleStruct := utils.Model2Article(articleModel)

	ctx.HTML(http.StatusOK, fmt.Sprintf("%s.html", articleStruct.Date), nil)
	return
}

func Init(ctx *gin.Context)()  {
	_ = os.RemoveAll("./cache")
	bin.FetchList()
	ctx.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}

func HtmlExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func outputPage(articleModel model.Article, c *gin.Context) {
	dateStr := time.Time(articleModel.Date).Format("20060102")
	cmd := model.Rdb.Get(model.Ctx, dateStr)
	pageStr := cmd.Val()
	if pageStr == "" {
		path := utils.GetAbsolutePath(articleModel)
		fileBytes,_ := ioutil.ReadFile(path)
		pageStr = string(fileBytes)
		model.Rdb.Set(model.Ctx, dateStr, pageStr, -1)
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, pageStr)
}
