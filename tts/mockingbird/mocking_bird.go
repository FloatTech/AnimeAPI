// Package mockingbird 拟声鸟
package mockingbird

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
)

const (
	dbpath          = "data/MockingBird/"
	cachePath       = dbpath + "cache/"
	azfile          = dbpath + "az.wav"
	wjfile          = dbpath + "wj.wav"
	ysgfile         = dbpath + "ysg.wav"
	baseURL         = "http://aaquatri.com/sound"
	synthesizersURL = baseURL + "/api/synthesizers/"
	synthesizeURL   = baseURL + "/api/synthesize"
	modeName        = "拟声鸟"
)

var (
	vocoderList      = []string{"WaveRNN", "HifiGAN"}
	mockingbirdModes = map[int]string{0: "阿梓", 1: "文静", 2: "药水哥"}
	exampleFileMap   = map[int]string{0: azfile, 1: wjfile, 2: ysgfile}
)

// MockingBirdTTS 类
type MockingBirdTTS struct {
	vocoder         int
	synt            int
	name            string
	exampleFileName string
}

// String 服务名
func (tts *MockingBirdTTS) String() string {
	return modeName + tts.name
}

func NewMockingBirdTTS(synt int) (*MockingBirdTTS, error) {
	if synt < 0 || synt < 1 {
		synt = 0
	}
	switch synt {
	case 0:
		_, err := file.GetLazyData(azfile, true)
		if err != nil {
			return nil, err
		}
	case 1:
		_, err := file.GetLazyData(wjfile, true)
		if err != nil {
			return nil, err
		}
	case 2:
		_, err := file.GetLazyData(ysgfile, true)
		if err != nil {
			return nil, err
		}
	}
	return &MockingBirdTTS{1, synt, mockingbirdModes[synt], exampleFileMap[synt]}, nil
}

// Speak 返回音频本地路径
func (tts *MockingBirdTTS) Speak(uid int64, text func() string) (fileName string, err error) {
	// 异步
	rch := make(chan string, 1)
	sch := make(chan string, 1)
	// 获得回复
	go func() {
		rch <- text()
	}()
	// 拟声器生成音频
	go func() {
		var spth string
		spth, err = tts.getSyntPath()
		sch <- spth
	}()
	fileName, err = tts.getWav(<-rch, <-sch, vocoderList[tts.vocoder], uid)
	if err != nil {
		return
	}
	fileName = "file:///" + file.BOTPATH + "/" + cachePath + fileName
	// 回复
	return
}

func (tts *MockingBirdTTS) getSyntPath() (syntPath string, err error) {
	data, err := web.RequestDataWith(web.NewDefaultClient(), synthesizersURL, "GET", "", "")
	if err != nil {
		return
	}
	syntPath = gjson.Get(binary.BytesToString(data), fmt.Sprintf("%d.path", tts.synt)).String()
	return
}

func (tts *MockingBirdTTS) getWav(text, syntPath, vocoder string, uid int64) (fileName string, err error) {
	if syntPath == "" {
		err = errors.New("nil syntPath")
		return
	}
	fileName = strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + "_mockingbird.wav"
	b := binary.SelectWriter()
	defer binary.PutWriter(b)
	w := multipart.NewWriter(b)
	// Add your file
	f, err := os.Open(tts.exampleFileName)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", tts.exampleFileName)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	if fw, err = w.CreateFormField("text"); err != nil {
		return
	}
	if _, err = fw.Write(binary.StringToBytes(text)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("synt_path"); err != nil {
		return
	}
	if _, err = fw.Write(binary.StringToBytes(syntPath)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("vocoder"); err != nil {
		return
	}
	if _, err = fw.Write(binary.StringToBytes(vocoder)); err != nil {
		return
	}
	_ = w.Close()
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", synthesizeURL, bytes.NewReader(b.Bytes()))
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = errors.New("bad status:" + res.Status)
		return
	}
	defer res.Body.Close()
	fo, err := os.Create(cachePath + fileName)
	if err != nil {
		return
	}
	defer fo.Close()
	_, err = io.Copy(fo, res.Body)
	if err != nil {
		return
	}
	return
}
