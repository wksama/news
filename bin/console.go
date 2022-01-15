package bin

import (
	"fmt"
	"news/cacher"
	"news/service/db"
	"news/service/sitemap"
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
	if articleModel.DateStr() != "" {
		dbDriver.Store(&articleModel)
		cacheDriver := cacher.New()
		cacheDriver.Store(&articleModel)
	}
	if group != nil {
		group.Done()
	}
}

func Cache() {
	database := new(db.Database)
	articleModels := database.List()
	cacheDriver := cacher.New()
	sitemaper := sitemap.New()
	var list []utils.ListItem
	for _, articleModel := range articleModels {
		item := utils.Model2ListItem(articleModel)
		list = append(list, item)
		cacheDriver.Store(&articleModel)
		sitemaper.Add(sitemap.ListItem2Link(item))
	}
	sitemaper.Save()
	cacheDriver.List()

	_ = os.WriteFile(utils.AbsolutPath("/cache/404.html"), []byte(tpl.RenderNotFoundPage()), 0777)
}
