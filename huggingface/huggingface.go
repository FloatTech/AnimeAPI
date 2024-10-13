// Package huggingface ai界的github, 常用参数
package huggingface

import (
	"bytes"
	"encoding/json"

	"github.com/FloatTech/floatbox/web"
)

const (
	// HuggingfaceSpaceHTTPS huggingface space https api
	HuggingfaceSpaceHTTPS = "https://hf.space"
	// Embed huggingface space api embed
	Embed = HuggingfaceSpaceHTTPS + "/embed"
	// HTTPSPushPath 推送队列
	HTTPSPushPath = Embed + "/%v/api/queue/push/"
	// HTTPSStatusPath 状态队列
	HTTPSStatusPath = Embed + "/%v/api/queue/status/"
	// HuggingfaceSpaceWss huggingface space wss api
	HuggingfaceSpaceWss = "wss://spaces.huggingface.tech"
	// WssJoinPath 推送队列2
	WssJoinPath = HuggingfaceSpaceWss + "/%v/queue/join"
	// HTTPSPredictPath 推送队列3
	HTTPSPredictPath = Embed + "/%v/api/predict/"
	// DefaultAction 默认动作
	DefaultAction = "predict"
	// CompleteStatus 完成状态
	CompleteStatus = "COMPLETE"
	// WssCompleteStatus 完成状态2
	WssCompleteStatus = "process_completed"
	// TimeoutMax 超时时间
	TimeoutMax = 300
)

// PushRequest 推送默认请求
type PushRequest struct {
	Action      string        `json:"action,omitempty"`
	FnIndex     int           `json:"fn_index"`
	Data        []interface{} `json:"data"`
	SessionHash string        `json:"session_hash"`
}

// PushResponse 推送默认响应
type PushResponse struct {
	Hash          string `json:"hash"`
	QueuePosition int    `json:"queue_position"`
}

// StatusRequest 状态默认请求
type StatusRequest struct {
	Hash string `json:"hash"`
}

// StatusResponse 状态默认响应
type StatusResponse struct {
	Status string `json:"status"`
	Data   struct {
		Data            []interface{} `json:"data"`
		Duration        float64       `json:"duration"`
		AverageDuration float64       `json:"average_duration"`
	}
}

// Push 推送请求
func Push(pushURL string, pushReq *PushRequest) (pushRes PushResponse, err error) {
	b, err := json.Marshal(pushReq)
	if err != nil {
		return
	}
	data, err := web.PostData(pushURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &pushRes)
	return
}

// Status 状态请求
func Status(statusURL string, statusReq *StatusRequest) (data []byte, err error) {
	b, err := json.Marshal(statusReq)
	if err != nil {
		return
	}
	data, err = web.PostData(statusURL, "application/json", bytes.NewReader(b))
	return
}
