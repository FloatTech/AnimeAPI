package classify

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const head = "https://sayuri.fumiama.top/dice?class=9&url="

var (
	comments = []string{
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
)

// Classify 图片打分 返回值：class dhash comment, data
func Classify(targetURL string, isNoNeedImg bool) (int, string, string, []byte) {
	if targetURL[0] != '&' {
		targetURL = url.QueryEscape(targetURL)
	}

	u := head + targetURL
	if isNoNeedImg {
		u += "&noimg=true"
	}
	resp, err := http.Get(u)

	if err != nil {
		log.Warnf("[AI打分] %v", err)
		return 0, "", "", nil
	}

	if isNoNeedImg {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warnf("[AI打分] %v", err)
			return 0, "", "", nil
		}
		dhash := gjson.GetBytes(data, "img").String()
		class := int(gjson.GetBytes(data, "class").Int())
		return class, dhash, comments[class], nil
	}

	class, err := strconv.Atoi(resp.Header.Get("Class"))
	dhash := resp.Header.Get("DHash")
	if err != nil {
		log.Warnf("[AI打分] %v", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Warnf("[AI打分] %v", err)
	}
	return class, dhash, comments[class], data
}
