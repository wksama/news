package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"html/template"
	"os"
	"penti/model"
	"strconv"
	"time"
)

const TPL = "./templates/article.gohtml"

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

func RenderHtml(data interface{}) *bytes.Buffer {
	tpl, _ := template.ParseFiles(TPL)
	var buf = new(bytes.Buffer)
	err := tpl.Execute(buf, data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return buf
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

func GetSaveDir(article model.Article) string {
	articleDate := time.Time(article.Date)
	dir := fmt.Sprintf("./cache/html/%d/%d/%d", articleDate.Year(), articleDate.Month(), articleDate.Day())
	os.MkdirAll(dir, 0777)

	return dir
}

func GetSaveName(article model.Article) string {
	articleDate := time.Time(article.Date).Format("20060102")

	return fmt.Sprintf("%s.html", articleDate)
}

func GetAbsolutePath(article model.Article) string {
	return fmt.Sprintf("%s/%s", GetSaveDir(article), GetSaveName(article))
}

func DateToFloat64(date datatypes.Date) float64 {
	float, err := strconv.ParseFloat(time.Time(date).Format("20060102"), 64)
	if err != nil {
		return 0
	}
	return float
}
