package model

import (
	"log"
	"news/resources"
)

func Init() {
	err := resources.Db.AutoMigrate(Article{})
	if err != nil {
		log.Fatalln("Failed to migrate model: ", err.Error())
	}
	//if viper.GetString("app.database") == "sqlite3" {
	//	var count int64
	//	resources.Db.Model(Article{}).Count(&count)
	//	mysqlDsn := viper.GetString("mysql.dsn")
	//	if count == 0 && viper.GetString("mysql.dsn") != "" {
	//		mysqlConn, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	//		if err != nil {
	//			panic(err)
	//		}
	//		var articles []Article
	//		mysqlConn.Find(&articles)
	//		for _, article := range articles {
	//			resources.Db.Create(article)
	//		}
	//
	//		mysqlDb, _ := mysqlConn.DB()
	//		mysqlDb.Close()
	//	}
	//}

}
