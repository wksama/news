package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"penti/model"
	"time"
)

const ARTICLE_TPL = "./templates/article.gohtml"

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

func GetPageContent(articleModel model.Article) string {
	dateStr := time.Time(articleModel.Date).Format("20060102")
	cmd := model.Rdb.Get(model.Ctx, dateStr)
	pageStr := cmd.Val()
	if pageStr == "" {
		path := GetAbsolutePath(articleModel)
		fileBytes,_ := ioutil.ReadFile(path)
		pageStr = string(fileBytes)
		model.Rdb.Set(model.Ctx, dateStr, pageStr, -1)
	}

	return pageStr
}

func RenderHtml(data interface{}) *bytes.Buffer {
	tpl, _ := template.ParseFiles(ARTICLE_TPL)
	var buf = new(bytes.Buffer)
	err := tpl.Execute(buf, data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return buf
}