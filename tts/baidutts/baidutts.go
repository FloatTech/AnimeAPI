// Package baidutts 百度文字转语音
package baidutts

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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
	modeName     = "百度"
)

var (
	BaiduttsModes = map[int]string{0: "女声", 1: "男声", 3: "度逍遥", 4: "度丫丫"}
)

// BaiduTTS 百度类
type BaiduTTS struct {
	per  int
	name string
}

// String 服务名
func (tts *BaiduTTS) String() string {
	return modeName + tts.name
}

// NewBaiduTTS 新的百度语音
func NewBaiduTTS(per int) *BaiduTTS {
	switch per {
	case 0, 1, 3, 4:
		return &BaiduTTS{per, BaiduttsModes[per]}
	default:
		return &BaiduTTS{4, BaiduttsModes[4]}
	}
}

// Speak 返回音频本地路径
func (tts *BaiduTTS) Speak(uid int64, text func() string) string {
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
	fileName := getWav(<-rch, <-tch, 5, tts.per, 5, 5, uid)
	// 回复
	return "file:///" + file.BOTPATH + "/" + cachePath + fileName
}

func getToken() (accessToken string) {
	data, err := web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(tokenURL, grantType, clientID, clientSecret), "GET", "", ua)
	if err != nil {
		log.Errorln("[baidutts]:", err)
	}
	accessToken = gjson.Get(binary.BytesToString(data), "access_token").String()
	return
}

func getWav(tex, tok string, vol, per, spd, pit int, uid int64) (fileName string) {
	fileName = strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + "_baidu.wav"

	cuid := fmt.Sprintf("%x", md5.Sum(binary.StringToBytes(tok)))
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
