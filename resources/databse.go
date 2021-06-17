package resources

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var Db *gorm.DB

func databaseInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/news?charset=utf8mb4&parseTime=True&loc=Local", viper.Get("mysql.user"), viper.GetString("mysql.password"), viper.Get("mysql.host"), viper.Get("mysql.port"))
	dbConfig := &gorm.Config{}
	if !viper.GetBool("app.debug") {
		dbConfig.Logger = logger.Default.LogMode(logger.Silent)
	}
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), dbConfig)
	if err != nil {
		log.Fatalln("Failed to connect to database!")
	}
}
