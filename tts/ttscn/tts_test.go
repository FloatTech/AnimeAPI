package ttscn

import "testing"

func TestTTS(t *testing.T) {
	tts, err := NewTTSCN("中文（普通话，简体）", "晓双（女 - 儿童）", KBRates[0])
	if err != nil {
		t.Fatal(err)
	}
	if tts.voice != "zh-CN-XiaoshuangNeural" {
		t.Fatal(tts.voice)
	}
	fn, err := tts.Speak(0, func() string { return "测试一下。" })
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fn)
}
