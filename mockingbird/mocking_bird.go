// Package mockingbird 拟声鸟
package mockingbird

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
)

const (
	prio            = 250
	dbpath          = "data/MockingBird/"
	cachePath       = dbpath + "cache/"
	dbfile          = dbpath + "降噪3.wav"
	baseURL         = "http://aaquatri.com/sound/"
	synthesizersURL = baseURL + "api/synthesizers/"
	synthesizeURL   = baseURL + "api/synthesize"
)

var (
	vocoderList = []string{"WaveRNN", "HifiGAN"}
)

// Speak 返回音频本地路径
func Speak(uid int64, text func() string) string {
	// 异步
	rch := make(chan string, 1)
	sch := make(chan string, 1)
	// 获得回复
	go func() {
		rch <- text()
	}()
	// 拟声器生成音频
	go func() {
		sch <- getSyntPath()
	}()
	fileName := getWav(<-rch, <-sch, vocoderList[1], uid)
	// 回复
	return "file:///" + file.BOTPATH + "/" + cachePath + fileName
}

func getSyntPath() (syntPath string) {
	data, err := web.ReqWith(synthesizersURL, "GET", "", "")
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	syntPath = gjson.Get(binary.BytesToString(data), "0.path").String()
	return
}

func getWav(text, syntPath, vocoder string, uid int64) (fileName string) {
	fileName = strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + ".wav"
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your file
	f, err := os.Open(dbfile)
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", dbfile)
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if _, err = io.Copy(fw, f); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if fw, err = w.CreateFormField("text"); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if _, err = fw.Write([]byte(text)); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if fw, err = w.CreateFormField("synt_path"); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if _, err = fw.Write([]byte(syntPath)); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if fw, err = w.CreateFormField("vocoder"); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	if _, err = fw.Write([]byte(vocoder)); err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	w.Close()
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", synthesizeURL, &b)
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		log.Errorf("[mockingbird]bad status: %s", res.Status)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	err = os.WriteFile(cachePath+fileName, data, 0666)
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	return
}
