// Package kimoi AI 匹配 kimoi 词库
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

// Response 回复结构
type Response struct {
	// Reply 文本
	Reply string `json:"reply"`
	// Confidence 置信度, 建议不要使用 < 0.5 或 > 0.95 的结果
	Confidence float64 `json:"confidence"`
}

// Chat 用户对 AI 说一句话
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
