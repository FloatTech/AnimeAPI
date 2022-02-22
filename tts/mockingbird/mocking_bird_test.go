package mockingbird

import (
	"fmt"
	"testing"
)

func TestNewMockingBirdTTS(t *testing.T) {
	tts := NewMockingBirdTTS(0)
	//fmt.Println(tts.Speak(int64(1), func() string {
	//	return "我爱你"
	//}))
	tts = NewMockingBirdTTS(1)
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
}
