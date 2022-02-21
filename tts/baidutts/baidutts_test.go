package baidutts

import (
	"fmt"
	"testing"
)

func TestNewBaiduTTS(t *testing.T) {
	tts := NewBaiduTTS(0)
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//tts = NewBaiduTTS(1)
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	//tts = NewBaiduTTS(3)
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	tts = NewBaiduTTS(4)
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
}
