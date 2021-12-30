package utils

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"html/template"
	"news/model"
	"strconv"
	"time"
)

type Article struct {
	FullTitle  string       `json:"full_title"`
	RealTitle  string       `json:"real_title"`
	Url        string       `json:"url"`
	Date       string       `json:"date"`
	Paragraphs []*Paragraph `json:"paragraphs"`
}

type ListItem struct {
	Score  float64
	Member string
}

type Paragraph struct {
	Subject string `json:"title"`
	Bodies  []Body `json:"body"`
}

type Body struct {
	Type    string        `json:"type"`
	Content template.HTML `json:"content"`
}

func Model2Article(article model.Article) Article {
	dateStr := time.Time(article.Date).Format("2006-01-02")
	var articleParagraphs []*Paragraph
	_ = json.Unmarshal(article.Paragraphs, &articleParagraphs)
	return Article{
		FullTitle:  article.FullTitle,
		RealTitle:  article.RealTitle,
		Url:        article.Url,
		Date:       dateStr,
		Paragraphs: articleParagraphs,
	}
}

func Model2Z(article model.Article) redis.Z {
	dateStr := article.DateStr()
	dateInt, _ := strconv.Atoi(dateStr)
	return redis.Z{
		Score:  float64(dateInt),
		Member: article.FullTitle,
	}
}

func Model2ListItem(article model.Article) ListItem {
	dateStr := article.DateStr()
	dateInt, _ := strconv.Atoi(dateStr)
	return ListItem{
		Score:  float64(dateInt),
		Member: article.FullTitle,
	}
}
