package cacher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"news/model"
	"news/service/tpl"
	"news/utils"
	"strconv"
)

type Redis struct {
	handler *redis.Client
	ctx     context.Context
}

func (r *Redis) connect() {
	r.ctx = context.Background()
	dsn := viper.GetString("redis.dsn")
	r.handler = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := r.handler.Ping(r.ctx).Result()
	if pong != "PONG" || err != nil {
		panic(err)
	}
}

func (r *Redis) Store(articleModel *model.Article) {
	score, err := strconv.ParseFloat(articleModel.DateStr(), 64)
	if err != nil {
		panic(err)
	}
	r.handler.ZAdd(r.ctx, "list", &redis.Z{
		Score:  score,
		Member: articleModel.FullTitle,
	})
	articleHtml := tpl.RenderArticle(utils.Model2Article(*articleModel))
	r.handler.Set(r.ctx, articleModel.DateStr(), articleHtml, 0)
}

func (r *Redis) Fetch(dateStr string) string {
	return r.handler.Get(r.ctx, dateStr).Val()
}

func (r *Redis) List() string {
	res := r.handler.ZRevRangeWithScores(r.ctx, "list", 0, -1).Val()
	var list []utils.ListItem
	for _, item := range res {
		list = append(list, utils.ListItem{
			Score:  item.Score,
			Member: item.Member.(string),
		})
	}

	return tpl.RenderList(list)
}
