package resources

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var Db *gorm.DB

func databaseInit() {
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
		Db, err = gorm.Open(sqlite.Open(dsn), dbConfig)
	}
	if err != nil {
		log.Fatalln("Failed to connect to database: ", err.Error())
	}
}
