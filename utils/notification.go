package utils

import (
	"net/http"
	"net/url"
)

func FangTang(title, desp string)  {
	var reqMsg = make(url.Values)
	reqMsg.Set("text", title)
	reqMsg.Set("desp", desp)
	http.PostForm("https://sc.ftqq.com/SCU108769T003ceb8a56c654f881454f51cccb2f045f2ea048ba11a.send", reqMsg)
}
