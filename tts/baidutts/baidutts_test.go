package baidutts

import (
	"fmt"
	"testing"
)

func TestNewBaiduTTS(t *testing.T) {
	clientID := "6ACWsLOg3b7OyGUKGfHZfbXa"
	clientSecret := "nA6WP1d05qBoUYqxplNAV1inf8IHGwj9"
	tts := NewBaiduTTS(0, clientID, clientSecret)
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
	tts = NewBaiduTTS(1, clientID, clientSecret)
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
	tts = NewBaiduTTS(3, clientID, clientSecret)
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
	tts = NewBaiduTTS(4, clientID, clientSecret)
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
}
