// Package genshin 原神vits语音api
package genshin

import (
	goBinary "encoding/binary"
	"fmt"
	"hash/crc64"
	"net/url"
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
	modeName  = "原神"
	cachePath = "data/gsvits/"
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

// Genshin 原神类
type Genshin struct {
	mode int
	name string
	code string
}

// String 服务名
func (tts *Genshin) String() string {
	return modeName + tts.name
}

// NewGenshin 新的原神语音
func NewGenshin(mode int, code string) *Genshin {
	return &Genshin{
		mode: mode,
		name: SoundList[mode],
		code: code,
	}
}

// Speak 返回音频 url
func (tts *Genshin) Speak(_ int64, text func() string) (fileName string, err error) {
	t := text()
	u := fmt.Sprintf(CNAPI, url.QueryEscape(tts.name), url.QueryEscape(
		// 将数字转文字
		re.ReplaceAllStringFunc(t, func(s string) string {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				logrus.Errorln("[tts]", err)
				return s
			}
			return numcn.EncodeFromFloat64(f)
		}),
	), tts.code)
	var b [8]byte
	goBinary.LittleEndian.PutUint64(b[:], uint64(tts.mode))
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	h.Write(b[:])
	_, _ = h.Write(binary.StringToBytes(u))
	n := fmt.Sprintf(cachePath+"%016x.ogg", h.Sum64())
	if file.IsExist(n) {
		fileName = "file:///" + file.BOTPATH + "/" + n
		return
	}
	data, err := web.GetData(u)
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
