package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"penti/model"
	"penti/resources"
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

func GetAbsolutePathByArticle(article model.Article) string {
	return fmt.Sprintf("%s/%s", GetSaveDir(article), GetSaveName(article))
}

func GetAbsolutePathByDateStr(dateStr string) string {
	date, _ := time.Parse("20060102", dateStr)
	dir := fmt.Sprintf("./cache/html/%d/%d/%d", date.Year(), date.Month(), date.Day())
	os.MkdirAll(dir, 0777)

	return fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%s.html", dateStr))
}

func GetPageContentByDateStr(dateStr string) (pageStr string) {
	cmd := resources.RC.Get(resources.Ctx, dateStr)
	pageStr = cmd.Val()
	if pageStr == "" {
		path := GetAbsolutePathByDateStr(dateStr)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return ""
		} else {
			fileBytes, _ := ioutil.ReadFile(path)
			pageStr = string(fileBytes)
		}
	}

	return
}

func RenderHtml(data Article) *bytes.Buffer {
	tpl, _ := template.ParseFiles(ARTICLE_TPL)
	var buf = new(bytes.Buffer)
	err := tpl.Execute(buf, data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return buf
}

func CacheFile(path string, content []byte) bool {
	err := ioutil.WriteFile(path, content, 0777)
	return err == nil
}
