package model

import (
	"log"
	"penti/resources"
)

func Init() {
	err := resources.Db.AutoMigrate(Article{})
	if err != nil {
		log.Fatalln("Failed to migrate model")
	}
}
