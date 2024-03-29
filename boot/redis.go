package boot

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

var RC *redis.Client
var Ctx = context.Background()

func initRedis() {
	if viper.GetString("app.cacher") == "redis" {
		dsn := viper.GetString("redis.dsn")
		RC = redis.NewClient(&redis.Options{
			Addr:     dsn,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		pong, err := RC.Ping(Ctx).Result()
		if pong != "PONG" {
			if err != nil {
				fmt.Println(err)
			}
			log.Fatalln("Failed to connect to redis server")
		}
	}
}
