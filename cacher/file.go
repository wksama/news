package cacher

import (
	"io/ioutil"
	"news/model"
	"news/utils"
	"os"
)

type File struct {
}

func (f File) Store(articleModel *model.Article) {
	path := utils.GetAbsolutePathByDateStr(articleModel.DateStr())
	if _, err := os.Stat(path); err == nil {
		return
	}
	articleStruct := utils.Model2Article(*articleModel)
	html := utils.RenderArticle(articleStruct)
	err := ioutil.WriteFile(path, []byte(html), 0777)
	if err != nil {
		panic("File store error")
	}
}

func (f File) Fetch(dateStr string) string {
	path := utils.GetAbsolutePathByDateStr(dateStr)
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic("File fetch error")
	}

	return string(contentBytes)
}

func (f File) List() string {
	path := utils.AbsolutDir("/cache/index.html")
	if _, err := os.Stat(path); err != nil {
		return ""
	}

	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic("File fetch list error")
	}

	return string(contentBytes)
}
