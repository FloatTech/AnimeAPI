package mockingbird

import (
	"fmt"
	"testing"
)

func TestNewMockingBirdTTS(t *testing.T) {
	tts, err := NewMockingBirdTTS(0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
	tts, err = NewMockingBirdTTS(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tts.Speak(int64(1), func() string {
		return "我爱你"
	}))
}
