package admin

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"html/template"
	"net/http"
	"news/bin"
	"news/boot"
	"news/model"
	"news/service/spider"
	"news/utils"
	"os"
	"strings"
	"time"
)

func CacheArticle(ctx *gin.Context) {
	var page = 1
	var size = 20
	for true {
		var offset = (page - 1) * size
		var articleModels []model.Article
		boot.Db.Offset(offset).Limit(size).Order("date DESC").Find(&articleModels)
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

func Redis2DB(ctx *gin.Context) {
	res := boot.RC.ZRangeWithScores(ctx, "articleList", 0, 100000)
	list := res.Val()
	for _, item := range list {
		dateStr := utils.Float64ToString(item.Score)
		cmd := boot.RC.Get(boot.Ctx, dateStr)
		pageStr := cmd.Val()
		pageReader, _ := goquery.NewDocumentFromReader(strings.NewReader(pageStr))
		articleModel := new(model.Article)
		dateTime, _ := time.Parse("20060102", dateStr)
		articleModel.Date = datatypes.Date(dateTime)
		articleModel.FullTitle = item.Member.(string)

		titleArr := strings.Split(articleModel.FullTitle, "】")
		articleModel.RealTitle = titleArr[1]
		articleModel.Url = pageReader.Find(".weui-footer__link").AttrOr("href", "")

		articleNode := pageReader.Find(".weui-article > section")
		paragraphLen := articleNode.Find("h2").Length()
		var paragraphs []*spider.Paragraph
		for i := 0; i < paragraphLen; i++ {
			var paragraph = new(spider.Paragraph)
			paragraph.Title = articleNode.Find("h2").Eq(i).Text()
			sectionNode := articleNode.Find("section").Eq(i)

			sectionNode.Find("p").Each(func(i int, selection *goquery.Selection) {
				var body spider.Body
				if selection.HasClass("img-div") {
					body.Type = "img"
					src := selection.Find("img").AttrOr("src", "")
					body.Content = template.HTML(src)
				} else {
					body.Type = "text"
					body.Content = template.HTML(selection.Text())
				}
				paragraph.Bodies = append(paragraph.Bodies, body)
			})
			paragraphs = append(paragraphs, paragraph)
		}
		paragraphsBytes, _ := json.Marshal(paragraphs)
		articleModel.Paragraphs = paragraphsBytes
		bin.InsertIntoDb(articleModel)
	}
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
	bin.FetchFlow(url, nil)
	ctx.JSON(http.StatusOK, gin.H{
		"Success": true,
	})
}
