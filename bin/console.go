package bin

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"news/cacher"
	"news/service/db"
	"news/service/spider"
	"news/utils"
	"os"
	"sync"
	"time"
)

func FetchFirstPage() {
	fmt.Println("fetching list")
	s := spider.New()

	urlList := s.FetchPageList()

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(urlList))
	for _, url := range urlList {
		go FetchAndStore(url, waitGroup)
	}
	waitGroup.Wait()
}

func FetchLatestArticle() {
	dateStr := time.Now().Format("20060102")

	FetchArticleByDateStr(dateStr)
}

func FetchArticleByDateStr(dateStr string) {
	cacheDriver := cacher.New()
	html := cacheDriver.Fetch(dateStr)
	if html == "" {
		s := spider.New()
		url, articleDateStr := s.FetchLatestArticleUrl()
		if articleDateStr != "" {
			FetchAndStore(url, nil)
		}
	}
}

func FetchAndStore(url string, group *sync.WaitGroup) {
	dbDriver := new(db.Database)
	s := spider.New()
	articleModel := s.FetchArticle(url)
	dbDriver.Store(&articleModel)
	cacheDriver := cacher.New()
	cacheDriver.Store(&articleModel)
	if group != nil {
		group.Done()
	}
}

func Cache() {
	database := new(db.Database)
	articleModels := database.List()
	cacheDriver := cacher.New()
	var list []redis.Z
	for _, articleModel := range articleModels {
		list = append(list, utils.Model2Z(articleModel))
		cacheDriver.Store(&articleModel)
	}
	html := utils.RenderList(list)
	err := os.WriteFile(utils.AbsolutDir("/cache/index.html"), []byte(html), 0777)
	if err != nil {
		panic(err)
	}
}

//func FetchFlow(url string, waitGroup *sync.WaitGroup) (articleModel model.Article) {
//	log.Printf("正在爬取文章：%s", url)
//	articleModel = fetchArticleByUrl(url)
//	if articleModel.FullTitle != "" {
//		log.Printf("文章标题：{%s}", articleModel.FullTitle)
//		InsertIntoDb(&articleModel)
//		if articleModel.ID != 0 {
//			CacheFlow(articleModel)
//		} else {
//			log.Println(articleModel.FullTitle + "写入数据库失败")
//		}
//	} else {
//		color.Red("文章不存在")
//	}
//	if waitGroup != nil {
//		waitGroup.Done()
//	}
//
//	return
//}
//
//func fetchArticleByUrl(url string) model.Article {
//	s := spider.New()
//	return s.FetchArticle(url)
//}

//func CacheFlow(articleModel model.Article) {
//	log.Println("正在执行缓存流程")
//	insertIntoRedis(articleModel)
//	CacheTemplate(articleModel)
//	log.Println("执行缓存流程完成")
//}
//
//func InsertIntoDb(articleModel *model.Article) {
//	log.Println("正在写入数据库")
//	result := boot.Db.Create(&articleModel)
//	if result.Error != nil {
//		log.Printf("Insert into database error: %s", result.Error.Error())
//	}
//	log.Println("写入数据库完成")
//}

//func insertIntoRedis(articleModel model.Article) {
//	if boot.RC != nil {
//		log.Println("正在写入Redis")
//		boot.RC.ZAdd(boot.Ctx, "articleList", &redis.Z{
//			Score:  utils.DateToFloat64(articleModel.Date),
//			Member: articleModel.FullTitle,
//		})
//		articleBytes, _ := json.Marshal(articleModel)
//		boot.RC.Set(boot.Ctx, articleModel.DateStr(), string(articleBytes), 0)
//		log.Println("写入Redis完成")
//	}
//}
//
//func CacheTemplate(articleModel model.Article) {
//	log.Println("正在渲染模板")
//	articleStruct := utils.Model2Article(articleModel)
//	htmlStr := utils.RenderArticle(articleStruct)
//	file := fmt.Sprintf("%s/%s", utils.GetSaveDir(articleModel), utils.GetSaveName(articleModel))
//
//	_ = ioutil.WriteFile(file, []byte(htmlStr), 0777)
//	log.Println("渲染结果缓存至静态文件完成")
//}
