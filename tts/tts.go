// Package tts 文字转语音库
package tts

type TTS interface {
	// Speak 返回音频本地路径
	Speak(key int64, text func() string) string
}
