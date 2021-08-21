package classify

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const head = "https://sayuri.fumiama.top/dice?class=9&url="

var (
	datapath  string
	cachefile string
	lastvisit = time.Now().Unix()
	comments  = []string{
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

// Init 设置 datapath
func Init(dataPath string) {
	datapath = dataPath
	os.RemoveAll(datapath) // 清除缓存
	err := os.MkdirAll(datapath, 0755)
	if err != nil {
		panic(err)
	}
	cachefile = datapath + "cache"
}

// Flush 刷新时间戳
func Flush() {
	lastvisit = time.Now().Unix()
}

// Canvisit 可以访问
func CanVisit(delay int64) bool {
	if time.Now().Unix()-lastvisit > delay {
		Flush()
		return true
	}
	return false
}

// Classify 图片打分 返回值：class lastvisit dhash comment
func Classify(targeturl string, noimg bool) (int, int64, string, string) {
	lv := lastvisit
	if targeturl[0] != '&' {
		targeturl = url.QueryEscape(targeturl)
	}
	get_url := head + targeturl
	if noimg {
		get_url += "&noimg=true"
	}
	resp, err := http.Get(get_url)
	if err != nil {
		log.Warnf("[AI打分] %v", err)
		return 0, 0, "", ""
	} else {
		if noimg {
			data, err1 := ioutil.ReadAll(resp.Body)
			if err1 == nil {
				dhash := gjson.GetBytes(data, "img").String()
				class := int(gjson.GetBytes(data, "class").Int())
				return class, lv, dhash, comments[class]
			} else {
				log.Warnf("[AI打分] %v", err1)
				return 0, 0, "", ""
			}
		} else {
			class, err1 := strconv.Atoi(resp.Header.Get("Class"))
			dhash := resp.Header.Get("DHash")
			if err1 != nil {
				log.Warnf("[AI打分] %v", err1)
			}
			defer resp.Body.Close()
			// 写入文件
			data, _ := ioutil.ReadAll(resp.Body)
			f, _ := os.OpenFile(cachefile+strconv.FormatInt(lv, 10), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			defer f.Close()
			f.Write(data)
			return class, lv, dhash, comments[class]
		}
	}
}
