package cacher

import (
	"github.com/spf13/viper"
	"news/model"
)

type Cacher interface {
	connect()
	Store(articleModel *model.Article)
	Fetch(dateStr string) string
	List() string
}

func New() Cacher {
	var c Cacher
	cacheDriver := viper.Get("app.cacher")
	switch cacheDriver {
	case "redis":
		c = new(Redis)
	case "file":
		fallthrough
	default:
		c = new(File)
	}

	c.connect()
	return c
}
