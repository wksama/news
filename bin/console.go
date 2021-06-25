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
	log.Println("redis查询完成")
	if latest != dateStr {
		log.Println("正在爬取最新数据")
		s := spider.New()
		url, articleDateStr := s.FetchLatestArticleUrl()
		if articleDateStr == dateStr {
			log.Println("成功获取到最新文章链接")
			articleModel := FetchFlow(url)
			log.Println("成功获取到最新文章内容")
			dateStr := articleModel.DateStr()
			log.Println("更新最新文章日期")
			resources.RC.Set(resources.Ctx, "latest", dateStr, 0)
		}
	}
	color.Yellow("======处理完成：" + dateStr + "======")
}

func FetchFlow(url string) (articleModel model.Article) {
	log.Printf("666:q正在爬取文章：%s", url)
	articleModel = fetchArticleByUrl(url)
	log.Printf("文章标题：{%s}", articleModel.FullTitle)
	insertIntoDb(&articleModel)
	if articleModel.ID != 0 {
		color.Yellow("执行CacheFlow")
		CacheFlow(articleModel)
		color.Yellow("执行CacheFlow完成")

		color.Yellow("Bark通知")
		go utils.Bark(articleModel.RealTitle, articleModel.DateStr())
	} else {
		log.Println(articleModel.FullTitle + "写入数据库失败")
	}

	return
}

func fetchArticleByUrl(url string) model.Article {
	s := spider.New()
	return s.FetchArticle(url)
}

func CacheFlow(articleModel model.Article) {
	log.Println("正在执行缓存流程")
	insertIntoRedis(articleModel)
	CacheTemplate(articleModel)
	log.Println("执行缓存流程完成")
}

func insertIntoDb(articleModel *model.Article) {
	log.Println("正在写入数据库")
	result := resources.Db.Create(&articleModel)
	if result.Error != nil {
		log.Printf("Insert into database error: %s", result.Error.Error())
	}
	log.Println("写入数据库完成")
}

func insertIntoRedis(articleModel model.Article) {
	log.Println("正在写入Redis")
	resources.RC.ZAdd(resources.Ctx, "articleList", &redis.Z{
		Score:  utils.DateToFloat64(articleModel.Date),
		Member: articleModel.FullTitle,
	})
	log.Println("写入Redis完成")
}

func CacheTemplate(articleModel model.Article) {
	log.Println("正在渲染模板")
	articleStruct := utils.Model2Article(articleModel)
	htmlBuffer := utils.RenderHtml(articleStruct)
	log.Println("渲染模板完成")

	log.Println("正在将渲染结果写如Redis")
	resources.RC.Set(resources.Ctx, articleModel.DateStr(), htmlBuffer.String(), 0)
	file := fmt.Sprintf("%s/%s", utils.GetSaveDir(articleModel), utils.GetSaveName(articleModel))
	log.Println("渲染结果写如Redis完成")

	log.Println("正在将渲染结果缓存至静态文件")
	_ = ioutil.WriteFile(file, htmlBuffer.Bytes(), 0777)
	log.Println("渲染结果缓存至静态文件完成")
}
