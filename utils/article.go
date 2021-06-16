package utils

import (
	"encoding/json"
	"html/template"
	"news/model"
	"time"
)

type Article struct {
	FullTitle  string       `json:"full_title"`
	RealTitle  string       `json:"real_title"`
	Url        string       `json:"url"`
	Date       string       `json:"date"`
	Paragraphs []*Paragraph `json:"paragraphs"`
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
