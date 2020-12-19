package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/datatypes"
	"io/ioutil"
	"net/http"
	"os"
	"penti/model"
	"penti/utils"
	"time"
)

const PAGE_SIZE = 20000000

func List(ctx *gin.Context) {
	res := model.Rdb.ZRevRangeWithScores(ctx, "articleList",0, PAGE_SIZE)
	if len(res.Val()) == 0 {
		var modelList = make([]model.Article, 10)
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

	outputPage(utils.GetAbsolutePath(articleModel), ctx)
	return
}

func Today(ctx *gin.Context) {
	var articleModel model.Article
	model.Db.Where("date = ?", datatypes.Date(time.Now())).First(&articleModel)

	articleStruct := utils.Model2Article(articleModel)

	ctx.HTML(http.StatusOK, fmt.Sprintf("%s.html", articleStruct.Date), nil)
	return
}

func HtmlExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func outputPage(path string, c *gin.Context) {
	pageStr,_ := ioutil.ReadFile(path)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(pageStr))
}
