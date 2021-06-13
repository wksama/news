package resources

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

var RC *redis.Client
var Ctx = context.Background()

func redisInit() {
	RC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.Get("redis.host"), viper.Get("redis.port")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, _ := RC.Ping(Ctx).Result()
	if pong != "PONG" {
		log.Fatalln("Failed to connect to redis server")
	}
}
