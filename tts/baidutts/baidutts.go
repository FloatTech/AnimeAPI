// Package baidutts 百度文字转语音
package baidutts

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
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
func (tts *BaiduTTS) Speak(uid int64, text func() string) (fileName string, err error) {
	// 异步
	rch := make(chan string, 1)
	tch := make(chan string, 1)
	// 获得回复
	go func() {
		rch <- text()
	}()
	// 取到token
	go func() {
		var tok string
		tok, err = getToken()
		tch <- tok
	}()
	tok := <-tch
	if tok == "" {
		return
	}
	fileName, err = getWav(<-rch, tok, 5, tts.per, 5, 5, uid)
	if err != nil {
		return
	}
	// 回复
	return "file:///" + file.BOTPATH + "/" + cachePath + fileName, nil
}

func getToken() (accessToken string, err error) {
	data, err := web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(tokenURL, grantType, clientID, clientSecret), "GET", "", ua)
	if err != nil {
		return
	}
	accessToken = gjson.Get(binary.BytesToString(data), "access_token").String()
	return
}

func getWav(tex, tok string, vol, per, spd, pit int, uid int64) (fileName string, err error) {
	fileName = strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + "_baidu.wav"

	cuid := fmt.Sprintf("%x", md5.Sum(binary.StringToBytes(tok)))
	payload := strings.NewReader(fmt.Sprintf("tex=%s&lan=zh&ctp=1&vol=%d&per=%d&spd=%d&pit=%d&cuid=%s&tok=%s", tex, vol, per, spd, pit, cuid, tok))

	data, err := web.PostData(ttsURL, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return
	}
	if json.Valid(data) {
		err = errors.New(gjson.ParseBytes(data).Get("err_msg").String())
		return
	}
	err = os.WriteFile(cachePath+fileName, data, 0666)
	return
}
