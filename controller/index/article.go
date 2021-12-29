package index

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"news/cacher"
)

func List(ctx *gin.Context) {
	cacheDriver := cacher.New()
	listContent := cacheDriver.List()
	if listContent == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, listContent)
	return

	//pageQuery := ctx.DefaultQuery("page", "1")
	//pageInt, err := strconv.Atoi(pageQuery)
	//if err != nil {
	//	pageInt = 1
	//}
	//offset := (pageInt - 1) * ListPageSize
	//var list []redis.Z
	//if boot.RC != nil {
	//	res := boot.RC.ZRevRangeWithScores(ctx, "articleList", 0, ListPageSize)
	//	if len(res.Val()) == 0 {
	//		var modelList []model.Article
	//		boot.Db.Select("full_title, date").Order("date desc").Find(&modelList)
	//
	//		for _, item := range modelList {
	//			boot.RC.ZAdd(ctx, "articleList", &redis.Z{
	//				Score:  utils.DateToFloat64(item.Date),
	//				Member: item.FullTitle,
	//			})
	//		}
	//		res = boot.RC.ZRevRangeWithScores(ctx, "articleList", int64(offset), ListPageSize)
	//	}
	//	list = res.Val()
	//}
	//ctx.HTML(http.StatusOK, "list.gohtml", list)
}

func Item(ctx *gin.Context) {
	cacheDriver := cacher.New()
	dateStr := ctx.Param("date")

	pageContent := cacheDriver.Fetch(dateStr)
	if pageContent == "" {
		ctx.HTML(http.StatusNotFound, "404.gohtml", nil)
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, pageContent)
	return
}
