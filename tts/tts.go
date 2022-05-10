// Package tts 文字转语音库
package tts

type TTS interface {
	// Speak 返回音频本地路径
	Speak(key int64, text func() string) (fileName string, err error)
	// String 获得实际使用的回复服务名
	String() string
}
