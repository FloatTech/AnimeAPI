// Package lolimi https://api.lolimi.cn/
package lolimi

import (
	goBinary "encoding/binary"
	"fmt"
	"hash/crc64"
	"os"
	"regexp"
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/pkumza/numcn"
	"github.com/sirupsen/logrus"
)

const (
	modeName  = "桑帛云"
	cachePath = "data/lolimi/"
)

var (
	re = regexp.MustCompile(`(\-|\+)?\d+(\.\d+)?`)
)

func init() {
	// _ = os.RemoveAll(cachePath)
	err := os.MkdirAll(cachePath, 0755)
	if err != nil {
		panic(err)
	}
}

// Lolimi 桑帛云 API
type Lolimi struct {
	mode int
	name string
}

// String 服务名
func (tts *Lolimi) String() string {
	return modeName + tts.name
}

// NewLolimi 新的桑帛云语音
func NewLolimi(mode int) *Lolimi {
	return &Lolimi{
		mode: mode,
		name: SoundList[mode],
	}
}

// Speak 返回音频 url
func (tts *Lolimi) Speak(_ int64, text func() string) (fileName string, err error) {
	t := text()
	// 将数字转文字
	t = re.ReplaceAllStringFunc(t, func(s string) string {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			logrus.Errorln("[tts]", err)
			return s
		}
		return numcn.EncodeFromFloat64(f)
	})
	var (
		b      [8]byte
		data   []byte
		ttsURL string
	)
	goBinary.LittleEndian.PutUint64(b[:], uint64(tts.mode))
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	h.Write(b[:])
	ttsURL, err = TTS(tts.name, t)
	if err != nil {
		return
	}
	_, _ = h.Write(binary.StringToBytes(ttsURL))
	n := fmt.Sprintf(cachePath+"%016x.wav", h.Sum64())
	if file.IsExist(n) {
		fileName = "file:///" + file.BOTPATH + "/" + n
		return
	}
	data, err = web.GetData(ttsURL)
	if err != nil {
		return
	}
	err = os.WriteFile(n, data, 0644)
	if err != nil {
		return
	}
	fileName = "file:///" + file.BOTPATH + "/" + n
	return
}
