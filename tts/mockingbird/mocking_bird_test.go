package mockingbird

import (
	"fmt"
	"testing"
)

func TestNewMockingBirdTTS(t *testing.T) {
	fmt.Println(NewMockingBirdTTS("药水哥"))
	fmt.Println(NewMockingBirdTTS("阿梓"))
}

func TestGetSyntPath(t *testing.T) {
	mb := NewMockingBirdTTS("药水哥")
	fmt.Println(mb.getSyntPath())
	mb = NewMockingBirdTTS("阿梓")
	fmt.Println(mb.getSyntPath())
}

func TestSpeak(t *testing.T) {
	mb := NewMockingBirdTTS("药水哥")
	text := func() string {
		return "你好，能不能做我老婆"
	}
	fmt.Println(mb.Speak(int64(1), text))
}
