package tts

import (
	"fmt"
	"testing"
)

func TestTTS(t *testing.T) {
	tts := NewTTS("百度女声")
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//tts=NewTTS("百度男声")
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//tts=NewTTS("百度度逍遥")
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//tts=NewTTS("百度度丫丫")
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//tts=NewTTS("拟声鸟阿梓")
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//fmt.Println(NewTTS("拟声鸟药水哥"))
	tts = NewTTS("拟声鸟药水哥")
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
}
