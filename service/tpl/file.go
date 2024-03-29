package tpl

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"html/template"
	"news/utils"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const ArticleTpl = "/templates/article.gohtml"
const ListTpl = "/templates/list.gohtml"
const NotFoundTpl = "/templates/404.gohtml"

func GetAbsolutePathByDateStr(dateStr string) string {
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		panic(err)
	}
	dir := fmt.Sprintf("%s/cache/%d/%02d", utils.ProjectRoot(), date.Year(), int(date.Month()))
	_ = os.MkdirAll(dir, 0777)

	return fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%s.html", dateStr))
}

func RenderArticle(data utils.Article) string {
	return RenderHtml(utils.AbsolutPath(ArticleTpl), data)
}

func RenderList(data []utils.ListItem) string {
	return RenderHtml(utils.AbsolutPath(ListTpl), data)
}

func RenderNotFoundPage() string {
	return RenderHtml(utils.AbsolutPath(NotFoundTpl), nil)
}

func RenderHtml(tplPath string, data interface{}) string {
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(template.FuncMap{
		"path": func(dateStr string) string {
			return GetUrlByDateStr(dateStr)
		},
		"home": func() string {
			baseUrl := viper.GetString("app.baseUrl")
			if baseUrl == "" {
				baseUrl = "/"
			}
			return baseUrl
		},
	}).ParseFiles(tplPath)
	if err != nil {
		panic(err)
	}
	var buf = new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	var minifiedBuffer = new(bytes.Buffer)
	err = m.Minify("text/html", minifiedBuffer, buf)
	if err != nil {
		color.Red("压缩HTML失败： ", err.Error())
		return ""
	}

	return minifiedBuffer.String()
}

func GetUrlByDateStr(dateStr string) string {
	if viper.GetBool("app.pages") {
		return fmt.Sprintf("%s/%s/%s/%s", viper.GetString("app.baseUrl"), dateStr[0:4], dateStr[4:6], dateStr)
	} else {
		return fmt.Sprintf("%s/date/%s", viper.GetString("app.baseUrl"), dateStr)
	}
}

//func GetPageContentByDateStr(dateStr string) (pageStr string) {
//	switch viper.GetString("app.read") {
//	case "file":
//		path := GetAbsolutePathByDateStr(dateStr)
//		_, err := os.Stat(path)
//		if os.IsNotExist(err) {
//			return ""
//		} else {
//			fileBytes, _ := ioutil.ReadFile(path)
//			pageStr = string(fileBytes)
//		}
//	case "redis":
//		cmd := boot.RC.Get(boot.Ctx, dateStr)
//		articleJson := cmd.Val()
//		var articleModel model.Article
//		_ = json.Unmarshal([]byte(articleJson), &articleModel)
//		article := Model2Article(articleModel)
//		pageStr = RenderArticle(article)
//	case "render":
//		var articleModel model.Article
//		boot.Db.Where("date = ?", Str2Date(dateStr)).First(&articleModel)
//		article := Model2Article(articleModel)
//		pageStr = RenderArticle(article)
//	}
//
//	return
//}
//func GetSaveDir(article model.Article) string {
//	articleDate := time.Time(article.Date)
//	dir := fmt.Sprintf("%s/cache/%d/%d/%d", ProjectRoot(), articleDate.Year(), articleDate.Month(), articleDate.Day())
//	os.MkdirAll(dir, 0777)
//
//	return dir
//}
//func GetSaveName(article model.Article) string {
//	articleDate := time.Time(article.Date).Format("20060102")
//
//	return fmt.Sprintf("%s.html", articleDate)
//}
