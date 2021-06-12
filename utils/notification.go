package utils

import (
	"fmt"
	"net/http"
	"net/url"
)

func FangTang(title, desp string)  {
	var reqMsg = make(url.Values)
	reqMsg.Set("text", title)
	reqMsg.Set("desp", desp)
	_, _ = http.PostForm("https://sc.ftqq.com/SCU108769T003ceb8a56c654f881454f51cccb2f045f2ea048ba11a.send", reqMsg)
}

func Bark(title, date string) {
	_, _ = http.DefaultClient.Get(fmt.Sprintf("https://api.day.app/i5DgrmdG5navWxCvgKT2Hc/%s/%s?sound=suspense&url=https://news.gomain.run/date/%s", date, title, date))
}
