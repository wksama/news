package cacher

import (
	"io/ioutil"
	"news/model"
	"news/service/tpl"
	"news/utils"
	"os"
)

type File struct {
}

func (f File) connect() {
	err := os.MkdirAll(tpl.AbsolutDir("/cache"), 0777)
	if err != nil {
		panic(err)
	}
}

func (f File) Store(articleModel *model.Article) {
	path := tpl.GetAbsolutePathByDateStr(articleModel.DateStr())
	if _, err := os.Stat(path); err == nil {
		return
	}
	articleStruct := utils.Model2Article(*articleModel)
	html := tpl.RenderArticle(articleStruct)
	err := ioutil.WriteFile(path, []byte(html), 0777)
	if err != nil {
		panic("File store error")
	}
}

func (f File) Fetch(dateStr string) string {
	path := tpl.GetAbsolutePathByDateStr(dateStr)
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
	path := tpl.AbsolutDir("/cache/index.html")
	if _, err := os.Stat(path); err != nil {
		return ""
	}

	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic("File fetch list error")
	}

	return string(contentBytes)
}
