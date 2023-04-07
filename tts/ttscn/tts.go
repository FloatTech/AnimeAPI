// Package ttscn https://www.text-to-speech.cn/
package ttscn

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
)

// TTS 类
type TTS struct {
	language    string // 中文（普通话，简体）
	voice       string // zh-CN-XiaoxiaoNeural
	role        int
	style       int
	rate        int
	pitch       int
	kbitrate    string // audio-16khz-32kbitrate-mono-mp3
	silence     string
	styledegree float64 // 1
}

// LangSpeakers ...
type LangSpeakers struct {
	ShortName []string // ShortName 用于 API
	LocalName []string // LocalName 用于显示
}

const (
	speakerlistapi = "https://www.text-to-speech.cn/getSpeekList.php"
	ttsapi         = "https://www.text-to-speech.cn/getSpeek.php"
)

//go:embed speakers.json
var embededspeakers []byte

var (
	Langs = func() (m map[string]*LangSpeakers) {
		data, err := web.GetData(speakerlistapi)
		if err != nil {
			_ = json.Unmarshal(embededspeakers, &m)
			return
		}
		_ = json.Unmarshal(data, &m)
		return
	}()
	// KBRates 质量
	KBRates = [...]string{
		"audio-16khz-32kbitrate-mono-mp3",
		"audio-16khz-128kbitrate-mono-mp3",
		"audio-24khz-160kbitrate-mono-mp3",
		"audio-48khz-192kbitrate-mono-mp3",
		"riff-16khz-16bit-mono-pcm",
		"riff-24khz-16bit-mono-pcm",
		"riff-48khz-16bit-mono-pcm",
	}
)

// String 服务名
func (tts *TTS) String() string {
	return tts.language + tts.voice
}

// NewMockingBirdTTS ...
func NewTTSCN(lang, speaker, kbrate string) (*TTS, error) {
	spks, ok := Langs[lang]
	if !ok {
		return nil, errors.New("no language named " + lang)
	}
	hasfound := false
	for i, s := range spks.LocalName {
		if s == speaker {
			speaker = spks.ShortName[i]
			hasfound = true
			break
		}
	}
	if !hasfound {
		for _, s := range spks.ShortName {
			if s == speaker {
				hasfound = true
				break
			}
		}
	}
	if !hasfound {
		return nil, errors.New("no speaker named " + speaker)
	}
	hasfound = false
	for _, s := range KBRates {
		if s == kbrate {
			hasfound = true
			break
		}
	}
	if !hasfound {
		return nil, errors.New("no kbrate named " + kbrate)
	}
	return &TTS{
		language:    lang,
		voice:       speaker,
		kbitrate:    kbrate,
		styledegree: 1,
	}, nil
}

type result struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Download string `json:"download"`
	Author   string `json:"author"`
	URL      string `json:"url"`
}

// Speak 返回音频本地路径
func (tts *TTS) Speak(_ int64, text func() string) (fileName string, err error) {
	q := binary.NewWriterF(func(w *binary.Writer) {
		w.WriteString("language=")
		w.WriteString(url.QueryEscape(tts.language))
		w.WriteString("&voice=")
		w.WriteString(url.QueryEscape(tts.voice))
		w.WriteString("&text=")
		w.WriteString(url.QueryEscape(text()))
		w.WriteString("&role=")
		w.WriteString(strconv.Itoa(tts.role))
		w.WriteString("&style=")
		w.WriteString(strconv.Itoa(tts.style))
		w.WriteString("&rate=")
		w.WriteString(strconv.Itoa(tts.rate))
		w.WriteString("&pitch=")
		w.WriteString(strconv.Itoa(tts.pitch))
		w.WriteString("&kbitrate=")
		w.WriteString(tts.kbitrate)
		w.WriteString("&silence=")
		w.WriteString(tts.silence)
		w.WriteString("&styledegree=")
		w.WriteString(strconv.FormatFloat(tts.styledegree, 'f', 2, 64))
	})
	println(string(q))
	data, err := web.RequestDataWithHeaders(
		web.NewTLS12Client(), ttsapi, "POST", func(r *http.Request) error {
			r.Header.Add("accept", "*/*")
			r.Header.Add("content-length", strconv.Itoa(len(q)))
			r.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
			r.Header.Add("origin", "https://www.text-to-speech.cn")
			r.Header.Add("referer", "https://www.text-to-speech.cn/")
			r.Header.Add("user-agent", web.RandUA())
			r.Header.Add("x-requested-with", "XMLHttpRequest")
			return nil
		}, bytes.NewReader(q))
	if err != nil {
		return
	}
	var re result
	err = json.Unmarshal(data, &re)
	if err != nil {
		return
	}
	if re.Code != 200 {
		err = errors.New(re.Msg)
		return
	}
	fileName = re.Download
	return
}
