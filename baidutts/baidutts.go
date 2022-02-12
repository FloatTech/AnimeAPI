// Package baidutts 百度文字转语音
package baidutts

import (
	"crypto/md5"
	"fmt"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	grantType    = "client_credentials"
	clientID     = "6ACWsLOg3b7OyGUKGfHZfbXa"
	clientSecret = "nA6WP1d05qBoUYqxplNAV1inf8IHGwj9"
	tokenURL     = "https://aip.baidubce.com/oauth/2.0/token?grant_type=%s&client_id=%s&client_secret=%s"
	dbpath       = "data/baidutts/"
	cachePath    = dbpath + "cache/"
	ttsURL       = "http://tsn.baidu.com/text2audio"
	ua           = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

// Speak 返回音频本地路径
func Speak(uid int64, per int, text func() string) string {
	// 异步
	rch := make(chan string, 1)
	tch := make(chan string, 1)
	// 获得回复
	go func() {
		rch <- text()
	}()
	// 取到token
	go func() {
		tch <- getToken()
	}()
	fileName := getWav(<-rch, <-tch, 5, per, 5, 5, uid)
	// 回复
	return "file:///" + file.BOTPATH + "/" + cachePath + fileName
}

func getToken() (accessToken string) {
	data, err := web.ReqWith(fmt.Sprintf(tokenURL, grantType, clientID, clientSecret), "GET", "", ua)
	if err != nil {
		log.Errorln("[baidutts]:", err)
	}
	accessToken = gjson.Get(helper.BytesToString(data), "access_token").String()
	return
}

func getWav(tex, tok string, vol, per, spd, pit int, uid int64) (fileName string) {
	fileName = strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + ".wav"

	cuid := fmt.Sprintf("%x", md5.Sum(helper.StringToBytes(tok)))
	payload := strings.NewReader(fmt.Sprintf("tex=%s&lan=zh&ctp=1&vol=%d&per=%d&spd=%d&pit=%d&cuid=%s&tok=%s", tex, vol, per, spd, pit, cuid, tok))

	client := &http.Client{}
	req, err := http.NewRequest("POST", ttsURL, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", ua)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)
	err = os.WriteFile(cachePath+fileName, data, 0666)
	if err != nil {
		log.Errorln("[baidutts]:", err)
	}
	return
}
