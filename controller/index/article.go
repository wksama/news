package index

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"news/model"
	"news/resources"
	"news/utils"
	"strconv"
)

const LIST_PAGE_SIZE = 20000000

func List(ctx *gin.Context) {
	pageQuery := ctx.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(pageQuery)
	if err != nil {
		pageInt = 1
	}
	offset := (pageInt - 1) * LIST_PAGE_SIZE
	res := resources.RC.ZRevRangeWithScores(ctx, "articleList", 0, LIST_PAGE_SIZE)
	if len(res.Val()) == 0 {
		var modelList []model.Article
		resources.Db.Select("full_title, date").Order("date desc").Find(&modelList)

		for _, item := range modelList {
			resources.RC.ZAdd(ctx, "articleList", &redis.Z{
				Score:  utils.DateToFloat64(item.Date),
				Member: item.FullTitle,
			})
		}
		res = resources.RC.ZRevRangeWithScores(ctx, "articleList", int64(offset), LIST_PAGE_SIZE)
	}
	list := res.Val()
	ctx.HTML(http.StatusOK, "list.gohtml", list)
}

func Item(ctx *gin.Context) {
	dateStr := ctx.Param("date")

	pageStr := utils.GetPageContentByDateStr(dateStr)
	if pageStr == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, pageStr)
	return
}
