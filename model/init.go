package model

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB
var Rdb *redis.Client

func InitDb() {
	dsn := "root:mysql2019@tcp(120.53.122.92:3333)/penti?charset=utf8mb4&parseTime=True&loc=Local"
	Db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	_ = Db.AutoMigrate(Article{})

	Rdb = redis.NewClient(&redis.Options{
		Addr:     "120.53.122.92:6377",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
