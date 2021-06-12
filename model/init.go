package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB
var Rdb *redis.Client
var Ctx = context.Background()

func InitDb() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config")   // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	dsn := fmt.Sprintf("root:mysql2019@tcp(%s:%s)/penti?charset=utf8mb4&parseTime=True&loc=Local", viper.Get("mysql.host"), viper.Get("mysql.port"))
	dbConfig := &gorm.Config{}
	if !viper.GetBool("app.debug") {
		dbConfig.Logger = logger.Default.LogMode(logger.Silent)
	}
	Db, _ = gorm.Open(mysql.Open(dsn), dbConfig)
	_ = Db.AutoMigrate(Article{})

	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.Get("redis.host"), viper.Get("redis.port")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
