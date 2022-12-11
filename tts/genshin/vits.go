package genshin

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/pkumza/numcn"
	"github.com/sirupsen/logrus"
)

const (
	modeName = "原神"
)

var (
	re = regexp.MustCompile(`(\-|\+)?\d+(\.\d+)?`)
)

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
func (tts *Genshin) Speak(uid int64, text func() string) (fileName string, err error) {
	fileName = fmt.Sprintf(cnapi, tts.mode, url.QueryEscape(
		// 将数字转文字
		re.ReplaceAllStringFunc(text(), func(s string) string {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				logrus.Errorln("[tts]", err)
				return s
			}
			return numcn.EncodeFromFloat64(f)
		}),
	), tts.code)
	return
}
