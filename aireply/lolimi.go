package aireply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

// LolimiAi Lolimi回复类
type LolimiAi struct {
	u string
	n string
	k string
	b []string
}

// LolimiMemoryAi Lolimi带记忆回复类
type LolimiMemoryAi struct {
	u string
	n string
	k string
	b []string
	l int
	m []lolimiMessage
}

// lolimiMessage 消息记忆记录
type lolimiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	lolimiURL = "https://apii.lolimi.cn"
	// MomoURL api地址
	MomoURL = lolimiURL + "/api/mmai/mm?key=%s&msg=%s"
	// MomoBotName ...
	MomoBotName = "沫沫"
	// JingfengURL api地址
	JingfengURL = lolimiURL + "/api/jjai/jj?key=%s&msg=%s"
	// JingfengBotName ...
	JingfengBotName = "婧枫"
	// GPT4oURL api地址
	GPT4oURL = lolimiURL + "/api/4o/gpt4o?key=%s&msg=%s"
	// GPT4oBotName ...
	// TODO 换个更好的名字
	GPT4oBotName = "GPT4o"

	// 带记忆 POST 请求专区

	// C4oURL api地址
	C4oURL = lolimiURL + "/api/c4o/c?key=%s"
	// C4oBotName ...
	// TODO 换个更好的名字
	C4oBotName = "GPT4o"
)

// NewLolimiAi ...
func NewLolimiAi(u, name string, key string, banwords ...string) *LolimiAi {
	return &LolimiAi{u: u, n: name, k: key, b: banwords}
}

// NewLolimiMemoryAi ...
func NewLolimiMemoryAi(u, name string, key string, limit int, banwords ...string) *LolimiMemoryAi {
	return &LolimiMemoryAi{u: u, n: name, k: key, l: limit, b: banwords, m: []lolimiMessage{}}
}

// LolimiAi

// String ...
func (l *LolimiAi) String() string {
	return l.n
}

// TalkPlain 取得回复消息
func (l *LolimiAi) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, l.n)
	u := fmt.Sprintf(l.u, url.QueryEscape(l.k), url.QueryEscape(msg))
	data, err := web.GetData(u)
	if err != nil {
		errMsg := err.Error()
		// Remove the key from error message
		errMsg = strings.ReplaceAll(errMsg, l.k, "********")
		return "ERROR: " + errMsg
	}
	replystr := gjson.Get(binary.BytesToString(data), "data.output").String()
	replystr = strings.ReplaceAll(replystr, "<img src=\"", "[CQ:image,file=")
	replystr = strings.ReplaceAll(replystr, "<br>", "\n")
	replystr = strings.ReplaceAll(replystr, "\" />", "]")
	textReply := strings.ReplaceAll(replystr, l.n, nickname)
	for _, w := range l.b {
		if strings.Contains(textReply, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	return textReply
}

// Talk 取得带 CQ 码的回复消息
func (l *LolimiAi) Talk(_ int64, msg, nickname string) string {
	return l.TalkPlain(0, msg, nickname)
}

// LolimiMemoryAi

// String ...
func (l *LolimiMemoryAi) String() string {
	return l.n
}

// TalkPlain 取得回复消息
func (l *LolimiMemoryAi) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, l.n)
	u := fmt.Sprintf(l.u, url.QueryEscape(l.k))
	json, err := json.Marshal(
		append(l.m,
			lolimiMessage{
				Role: "user", Content: msg,
			},
		),
	)
	if err != nil {
		//panic(err)
		return "ERROR: " + err.Error()
	}
	// TODO: 可能会返回
	// "请使用psot格式请求如有疑问进官方群"
	data, err := web.PostData(u, "application/json", bytes.NewReader(json))
	if err != nil {
		errMsg := err.Error()
		// Remove the key from error message
		errMsg = strings.ReplaceAll(errMsg, l.k, "********")
		return "ERROR: " + errMsg
	}
	replystr := binary.BytesToString(data)
	replystr = strings.ReplaceAll(replystr, "<img src=\"", "[CQ:image,file=")
	replystr = strings.ReplaceAll(replystr, "<br>", "\n")
	replystr = strings.ReplaceAll(replystr, "\" />", "]")
	textReply := strings.ReplaceAll(replystr, l.n, nickname)
	for _, w := range l.b {
		if strings.Contains(textReply, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	if len(l.m) >= l.l-1 {
		l.m = append(l.m[2:],
			lolimiMessage{
				Role: "user", Content: msg,
			},
			lolimiMessage{
				Role: "assistant", Content: textReply,
			},
		)
	} else {
		l.m = append(l.m,
			lolimiMessage{
				Role: "user", Content: msg,
			},
			lolimiMessage{
				Role: "assistant", Content: textReply,
			},
		)
	}
	return textReply
}

// Talk 取得带 CQ 码的回复消息
func (l *LolimiMemoryAi) Talk(_ int64, msg, nickname string) string {
	return l.TalkPlain(0, msg, nickname)
}
