package kimoi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	base14 "github.com/fumiama/go-base16384"
)

const key = "暻撈莬穔僿貶稙棯悟澸滰蓱咜唕母屬石褤汴儱榅璕婴㴅"

const api = "https://ninex.azurewebsites.net/api/chat?code="

type Response struct {
	Reply      string  `json:"reply"`
	Confidence float64 `json:"confidence"`
}

func Chat(msg string) (r Response, err error) {
	resp, err := http.Post(
		api+base64.URLEncoding.EncodeToString(base14.DecodeFromString(key)),
		"text/plain", bytes.NewBufferString(msg),
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&r)
	return
}
