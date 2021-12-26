package cacher

import (
	"news/model"
)

type Cacher interface {
	Store(articleModel *model.Article)
	Fetch(dateStr string) string
	List() string
}

func New() Cacher {
	//var d Cacher
	//cacheDriver := viper.Get("app.cache")
	//switch cacheDriver {
	//case "file":
	//	d = new(File)
	//case "redis":
	//	d =
	//}
	return new(File)
}
