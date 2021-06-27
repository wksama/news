package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	gorm.Model
	ID         uint           `gorm:"primarykey" json:"id"`
	FullTitle  string         `json:"full_title"`
	RealTitle  string         `json:"real_title"`
	Url        string         `gorm:"index:,unique;not null" json:"url"`
	Date       datatypes.Date `gorm:"index:,unique;not null" json:"date"`
	Paragraphs datatypes.JSON `json:"paragraphs"`
}

func (a *Article) DateStr() string {
	if a.Date != (datatypes.Date{}) {
		return time.Time(a.Date).Format("20060102")
	}

	return ""
}
