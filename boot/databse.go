package boot

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"news/model"
)

var Db *gorm.DB

func initDatabase() {
	dbConfig := &gorm.Config{}
	if !viper.GetBool("app.debug") {
		dbConfig.Logger = logger.Default.LogMode(logger.Silent)
	}
	database := viper.GetString("app.database")
	dsn := viper.GetString(fmt.Sprintf("%s.dsn", database))
	var err error
	switch database {
	case "mysql":
		Db, err = gorm.Open(mysql.Open(dsn), dbConfig)
	case "sqlite3":
		fallthrough
	default:
		Db, err = gorm.Open(sqlite.Open(dsn), dbConfig)
	}
	if err != nil {
		log.Fatalln("Failed to connect to database: ", err.Error())
	}

	err = Db.AutoMigrate(model.Article{})
	if err != nil {
		log.Fatalln("Failed to migrate model: ", err.Error())
	}
}
