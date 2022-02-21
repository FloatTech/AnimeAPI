// Package mockingbird 拟声鸟
package mockingbird

import (
	"bytes"
	"fmt"
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
	dbpath          = "data/MockingBird/"
	cachePath       = dbpath + "cache/"
	azfile          = dbpath + "az.wav"
	ysgfile         = dbpath + "ysg.wav"
	baseURL         = "http://aaquatri.com/sound/"
	synthesizersURL = baseURL + "api/synthesizers/"
	synthesizeURL   = baseURL + "api/synthesize"
)

var (
	syntNumber  int64
	vocoderList = []string{"WaveRNN", "HifiGAN"}
	modeMap     = func() (m map[string]*MockingBirdTTS) {
		setReplyMap := func(m map[string]*MockingBirdTTS, r *MockingBirdTTS) {
			m[r.name] = r
		}
		m = make(map[string]*MockingBirdTTS, 2)
		setReplyMap(m, &MockingBirdTTS{1, 0, "阿梓", azfile})
		setReplyMap(m, &MockingBirdTTS{1, 1, "药水哥", ysgfile})
		return
	}()
)

type MockingBirdTTS struct {
	vocoder         int
	synt            int
	name            string
	exampleFileName string
}

func NewMockingBirdTTS(mode string) *MockingBirdTTS {
	return modeMap[mode]
}

// Speak 返回音频本地路径
func (tts *MockingBirdTTS) Speak(uid int64, text func() string) string {
	// 异步
	rch := make(chan string, 1)
	sch := make(chan string, 1)
	// 获得回复
	go func() {
		rch <- text()
	}()
	// 拟声器生成音频
	go func() {
		sch <- tts.getSyntPath()
	}()
	fileName := tts.getWav(<-rch, <-sch, vocoderList[tts.vocoder], uid)
	// 回复
	return "file:///" + file.BOTPATH + "/" + cachePath + fileName
}

func (tts *MockingBirdTTS) getSyntPath() (syntPath string) {
	data, err := web.ReqWith(synthesizersURL, "GET", "", "")
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	syntNumber = gjson.Get(binary.BytesToString(data), "#").Int()
	syntPath = gjson.Get(binary.BytesToString(data), fmt.Sprintf("%d.path", tts.synt)).String()
	return
}

func (tts *MockingBirdTTS) getWav(text, syntPath, vocoder string, uid int64) (fileName string) {
	fileName = strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + "_mockingbird.wav"
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your file
	f, err := os.Open(tts.exampleFileName)
	if err != nil {
		log.Errorln("[mockingbird]:", err)
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", tts.exampleFileName)
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
