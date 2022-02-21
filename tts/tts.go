// Package tts 文字转语音库
package tts

import (
	"github.com/FloatTech/AnimeAPI/tts/baidutts"
	"github.com/FloatTech/AnimeAPI/tts/mockingbird"
)

var (
	modeMap = func() (m map[string]TTS) {
		setReplyMap := func(m map[string]TTS, r TTS) {
			m[r.String()] = r
		}
		m = make(map[string]TTS, 6)
		setReplyMap(m, &baidutts.BaiduTTS{Per: 0, Name: baidutts.BaiduttsModes[0]})
		setReplyMap(m, &baidutts.BaiduTTS{Per: 1, Name: baidutts.BaiduttsModes[1]})
		setReplyMap(m, &baidutts.BaiduTTS{Per: 3, Name: baidutts.BaiduttsModes[3]})
		setReplyMap(m, &baidutts.BaiduTTS{Per: 4, Name: baidutts.BaiduttsModes[4]})
		setReplyMap(m, &mockingbird.MockingBirdTTS{1, 0, "阿梓", mockingbird.Azfile})
		setReplyMap(m, &mockingbird.MockingBirdTTS{1, 1, "药水哥", mockingbird.Ysgfile})
		return
	}()
)

type TTS interface {
	// Speak 返回音频本地路径
	Speak(key int64, text func() string) string
	// String 获得实际使用的回复服务名
	String() string
}

// NewTTS智能回复简单工厂
func NewTTS(mode string) TTS {
	return modeMap[mode]
}
