package cacher

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"news/model"
	"news/service/tpl"
	"news/utils"
	"os"
	"sort"
	"strconv"
)

type File struct {
}

func (f File) connect() {
	err := os.MkdirAll(utils.AbsolutPath("/cache"), 0777)
	if err != nil {
		panic(err)
	}
}

func (f File) Store(articleModel *model.Article) {
	path := tpl.GetAbsolutePathByDateStr(articleModel.DateStr())
	//if _, err := os.Stat(path); err == nil {
	//	return
	//}
	listJsonFilePath := utils.AbsolutPath("/cache/list.json")
	var listJsonFileContent []byte
	if _, err := os.Stat(listJsonFilePath); err != nil {
		listJsonFile, err := os.Create(listJsonFilePath)
		if err != nil {
			panic(err)
		}
		defer listJsonFile.Close()
		listJsonFileContent, err = io.ReadAll(listJsonFile)
		if err != nil {
			panic(err)
		}
	} else {
		listJsonFileContent, err = ioutil.ReadFile(listJsonFilePath)
		if err != nil {
			panic(err)
		}
	}
	listJsonMap := make(map[string]string)
	_ = json.Unmarshal(listJsonFileContent, &listJsonMap)
	listJsonMap[articleModel.DateStr()] = articleModel.FullTitle
	listJsonFileContent, err := json.Marshal(listJsonMap)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(listJsonFilePath, listJsonFileContent, 0777)
	if err != nil {
		panic(err)
	}

	articleStruct := utils.Model2Article(*articleModel)
	html := tpl.RenderArticle(articleStruct)
	err = ioutil.WriteFile(path, []byte(html), 0777)
	if err != nil {
		panic("File store error")
	}
	err = os.Remove(utils.AbsolutPath("/cache/index.html"))
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
	listFileContent, err := ioutil.ReadFile(utils.AbsolutPath("/cache/index.html"))
	if !os.IsNotExist(err) {
		return string(listFileContent)
	}
	listJsonFileContent, err := ioutil.ReadFile(utils.AbsolutPath("/cache/list.json"))
	if !os.IsNotExist(err) {
		listJsonMap := make(map[string]string)
		err := json.Unmarshal(listJsonFileContent, &listJsonMap)
		if err != nil {
			panic(err)
		}
		var dateArr []string
		for dateStr, _ := range listJsonMap {
			dateArr = append(dateArr, dateStr)
		}
		sort.Slice(dateArr, func(i, j int) bool {
			return dateArr[i] > dateArr[j]
		})
		var list []utils.ListItem
		for _, dateStr := range dateArr {
			fullTitle := listJsonMap[dateStr]
			score, _ := strconv.ParseFloat(dateStr, 64)
			list = append(list, utils.ListItem{
				Score:  score,
				Member: fullTitle,
			})
		}
		html := tpl.RenderList(list)
		err = os.WriteFile(utils.AbsolutPath("/cache/index.html"), []byte(html), 0777)
		if err != nil {
			panic(err)
		}

		return html
	}

	return ""
}
