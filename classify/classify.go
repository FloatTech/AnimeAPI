package classify

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tidwall/gjson"
)

const head = "https://sayuri.fumiama.top/dice?class=9&url="

var Comments = [...]string{
	"[0]这啥啊",
	"[1]普通欸",
	"[2]有点可爱",
	"[3]不错哦",
	"[4]很棒",
	"[5]我好啦!",
	"[6]影响不好啦!",
	"[7]太涩啦，🐛了!",
	"[8]已经🐛不动啦...",
}

// Classify 图片打分
func Classify(targetURL string, isNoNeedImg bool) (class int, dhash string, data []byte, err error) {
	if targetURL[0] != '&' {
		targetURL = url.QueryEscape(targetURL)
	}

	u := head + targetURL
	if isNoNeedImg {
		u += "&noimg=true"
	}

	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if isNoNeedImg {
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		dhash = gjson.GetBytes(data, "img").String()
		class = int(gjson.GetBytes(data, "class").Int())
		return
	}

	class, err = strconv.Atoi(resp.Header.Get("Class"))
	dhash = resp.Header.Get("DHash")
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
