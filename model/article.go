package model

import (
	"gorm.io/datatypes"
	"time"
)

type Article struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  time.Time      `gorm:"index" json:"deleted_at"`
	FullTitle  string         `json:"full_title"`
	RealTitle  string         `json:"real_title"`
	Url        string         `gorm:"index:,unique" json:"url"`
	Date       datatypes.Date `gorm:"index:,unique,sort:desc,type: date" json:"date"`
	Paragraphs datatypes.JSON `json:"paragraphs"`
}

func (a *Article) DateStr() string {
	if a.Date != (datatypes.Date{}) {
		return time.Time(a.Date).Format("2016-01-02")
	}

	return ""
}
