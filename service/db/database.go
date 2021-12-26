package db

import (
	"log"
	"news/boot"
	"news/model"
	"news/utils"
)

type Database struct {
}

func (d Database) Store(articleModel *model.Article) {
	result := boot.Db.Create(articleModel)
	if result.Error != nil {
		log.Printf("Insert into database error: %s", result.Error.Error())
	}
}

func (d Database) Fetch(dateStr string) model.Article {
	var articleModel model.Article
	boot.Db.Where("date = ?", utils.Str2Date(dateStr)).First(&articleModel)

	return articleModel
}

func (d Database) List() []model.Article {
	var articleModels []model.Article
	boot.Db.Order("date DESC").Find(&articleModels)

	return articleModels
}
