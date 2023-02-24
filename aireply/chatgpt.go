package aireply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ChatGPT GPT回复类
type ChatGPT struct {
	u string
	k string
	b []string
}

// chatGPTResponseBody 响应体
type chatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

// chatGPTRequestBody 请求体
type chatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

const (
	// ChatGPTURL api地址
	ChatGPTURL = "https://api.openai.com/v1/"
)

// NewChatGPT ...
func NewChatGPT(u, key string, banwords ...string) *ChatGPT {
	return &ChatGPT{u: u, k: key, b: banwords}
}

// String ...
func (*ChatGPT) String() string {
	return "ChatGPT"
}

// Talk 取得带 CQ 码的回复消息
func (c *ChatGPT) Talk(_ int64, msg, nickname string) string {
	replystr, err := chat(msg, c.k, c.u)
	if err != nil {
		return err.Error()
	}
	for _, w := range c.b {
		if strings.Contains(replystr, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	return replystr
}

// TalkPlain 取得回复消息
func (c *ChatGPT) TalkPlain(_ int64, msg, nickname string) string {
	return c.Talk(0, msg, nickname)
}

func chat(msg string, apiKey string, url string) (string, error) {
	requestBody := chatGPTRequestBody{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        2048,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData := new(bytes.Buffer)
	err := json.NewEncoder(requestData).Encode(requestBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url+"completions", requestData)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var gptResponseBody chatGPTResponseBody
	err = json.NewDecoder(response.Body).Decode(&gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = fmt.Sprint(v["text"])
			break
		}
	}
	return reply, nil
}
