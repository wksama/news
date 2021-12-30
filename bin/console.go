package bin

import (
	"fmt"
	"github.com/spf13/viper"
	"news/cacher"
	"news/service/db"
	"news/service/spider"
	"news/service/tpl"
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
	var list []utils.ListItem
	for _, articleModel := range articleModels {
		list = append(list, utils.Model2ListItem(articleModel))
		cacheDriver.Store(&articleModel)
	}
	if viper.GetString("app.cacher") == "file" {
		html := tpl.RenderList(list)
		err := os.WriteFile(tpl.AbsolutDir("/cache/index.html"), []byte(html), 0777)
		if err != nil {
			panic(err)
		}
	}

	_ = os.WriteFile(tpl.AbsolutDir("/cache/404.html"), []byte(tpl.RenderNotFoundPage()), 0777)
}
