package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	gorm.Model
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	FullTitle  string         `json:"full_title"`
	RealTitle  string         `json:"real_title"`
	Url        string         `gorm:"index:,unique" json:"url"`
	Date       datatypes.Date `gorm:"index:,unique,sort:desc,type: date" json:"date"`
	Paragraphs datatypes.JSON `json:"paragraphs"`
}

func (a *Article) DateStr() string {
	if a.Date != (datatypes.Date{}) {
		return time.Time(a.Date).Format("20160102")
	}

	return ""
}
