package bin

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"news/model"
	"news/resources"
	"news/spider"
	"news/utils"
	"time"
)

func FetchFirstPage() {
	fmt.Println("fetching list")
	s := spider.New()

	var urlList []string
	urlList = append(urlList, s.FetchPageList()...)

	for index, url := range urlList {
		fmt.Println(fmt.Sprintf("index: %d", index))
		FetchFlow(url)
		fmt.Println("sleeping ...")
		time.Sleep(2 * time.Second)
	}
}

func FetchLatestArticle() {
	dateStr := time.Now().Format("20060102")
	color.Yellow("======正在处理：" + dateStr + "======")
	fmt.Println("Date: " + dateStr)
	latest, err := resources.RC.Get(resources.Ctx, "latest").Result()
	if err != nil {
		log.Fatalln("Redis查询错误")
	}
	if latest != dateStr {
		s := spider.New()
		url, articleDateStr := s.FetchLatestArticleUrl()
		if articleDateStr == dateStr {
			FetchFlow(url)
		}
	}
	color.Yellow("======处理完成：" + dateStr + "======")
}

func FetchFlow(url string) {
	log.Printf("正在爬取文章：%s", url)
	articleModel := fetchArticleByUrl(url)
	log.Printf("文章标题：{%s}", articleModel.FullTitle)
	insertIntoDb(&articleModel)
	if articleModel.ID != 0 {
		log.Println(articleModel.FullTitle + "写入数据库成功")
		CacheFlow(articleModel)
		go utils.Bark(articleModel.RealTitle, articleModel.DateStr())
	} else {
		log.Println(articleModel.FullTitle + "写入数据库失败")
	}
}

func fetchArticleByUrl(url string) model.Article {
	s := spider.New()
	return s.FetchArticle(url)
}

func CacheFlow(articleModel model.Article) {
	insertIntoRedis(articleModel)
	log.Println(articleModel.FullTitle + "写入数Redis成功")
	CacheTemplate(articleModel)
	log.Println(articleModel.FullTitle + "缓存模板成功")
}

func insertIntoDb(articleModel *model.Article) {
	result := resources.Db.Create(&articleModel)
	if result.Error != nil {
		log.Printf("Insert into database error: %s", result.Error.Error())
	}
}

func insertIntoRedis(articleModel model.Article) {
	dateStr := articleModel.DateStr()
	resources.RC.Set(resources.Ctx, "latest", dateStr, 0)
	resources.RC.ZAdd(resources.Ctx, "articleList", &redis.Z{
		Score:  utils.DateToFloat64(articleModel.Date),
		Member: articleModel.FullTitle,
	})
}

func CacheTemplate(articleModel model.Article) {
	articleStruct := utils.Model2Article(articleModel)
	htmlBuffer := utils.RenderHtml(articleStruct)
	resources.RC.Set(resources.Ctx, articleModel.DateStr(), htmlBuffer.String(), 0)
	file := fmt.Sprintf("%s/%s", utils.GetSaveDir(articleModel), utils.GetSaveName(articleModel))
	_ = ioutil.WriteFile(file, htmlBuffer.Bytes(), 0777)
}
