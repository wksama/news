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
	case "file":
		c = new(File)
	case "redis":
		c = new(Redis)
	}
	c.connect()
	return c
}
