package utils

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"net/url"
)

func FangTang(title, desp string) {
	var reqMsg = make(url.Values)
	reqMsg.Set("text", title)
	reqMsg.Set("desp", desp)
	_, _ = http.PostForm("https://sc.ftqq.com/SCU108769T003ceb8a56c654f881454f51cccb2f045f2ea048ba11a.send", reqMsg)
}

func Bark(title, date string) {
	barkUrl := fmt.Sprintf("https://api.day.app/i5DgrmdG5navWxCvgKT2Hc/%s/%s?sound=suspense&url=https://news.gomain.run/date/%s", date, title, date)
	log.Println(barkUrl)
	_, err := http.DefaultClient.Get(barkUrl)
	if err != nil {
		color.Red("Bark失败" + err.Error())
	}
}
