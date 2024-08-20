package aireply

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

// LolimiAi2 Lolimi回复类
type LolimiAi2 struct {
	u string
	n string
	k string
	b []string
}

// LolimiAi2Mem Lolimi带记忆回复类
type LolimiAi2Mem struct {
	u string
	n string
	k string
	b []string
	l int
	m []lolimi2Message
}

// lolimi2Message 消息记忆记录
type lolimi2Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	lolimi2URL = "https://apii.lolimi.cn"
	// Momo2URL api地址
	Momo2URL = lolimi2URL + "/api/mmai/mm?key=%s&msg=%v"
	// Momo2BotName ...
	Momo2BotName = "沫沫"
	// Jingfeng2URL api地址
	Jingfeng2URL = lolimi2URL + "/api/jjai/jj?key=%s&msg=%v"
	// Jingfeng2BotName ...
	Jingfeng2BotName = "婧枫"
	// GPT4oURL api地址
	GPT4oURL = lolimi2URL + "/api/4o/gpt4o?key=%s&msg=%v"
	// GPT4oBotName ...
	// TODO 换个更好的名字
	GPT4oBotName = "GPT4o"

	// 带记忆 POST 请求专区
	// C4oURL api地址
	C4oURL = lolimi2URL + "/api/c4o/c?key=%s"
	// C4oBotName ...
	// TODO 换个更好的名字
	C4oBotName = "GPT4o"
)

// NewLolimiAi2 ...
func NewLolimiAi2(u, name string, key string, banwords ...string) *LolimiAi2 {
	return &LolimiAi2{u: u, n: name, k: key, b: banwords}
}

// NewLolimiAi2Mem ...
func NewLolimiAi2Mem(u, name string, key string, limit int, banwords ...string) *LolimiAi2Mem {
	return &LolimiAi2Mem{u: u, n: name, k: key, l: limit, b: banwords, m: make([]lolimi2Message, limit)}
}

// String ...
func (l *LolimiAi2) String() string {
	return l.n
}

// String ...
func (l *LolimiAi2Mem) String() string {
	return l.n
}

// TalkPlain 取得回复消息
func (l *LolimiAi2) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, l.n)
	u := fmt.Sprintf(l.u, url.QueryEscape(l.k), url.QueryEscape(msg))
	data, err := web.GetData(u)
	if err != nil {
		return "ERROR: " + err.Error()
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

// TalkPlain 取得回复消息
func (l *LolimiAi2Mem) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, l.n)
	u := fmt.Sprintf(l.u, url.QueryEscape(l.k))
	json, err := json.Marshal(
		append(l.m,
			lolimi2Message{
				Role: "user", Content: msg,
			},
		),
	)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	data, err := web.PostData(u, "application/json", strings.NewReader(binary.BytesToString(json)))
	if err != nil {
		return "ERROR: " + err.Error()
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
			lolimi2Message{
				Role: "user", Content: msg,
			},
			lolimi2Message{
				Role: "assistant", Content: textReply,
			},
		)
	} else {
		l.m = append(l.m,
			lolimi2Message{
				Role: "user", Content: msg,
			},
			lolimi2Message{
				Role: "assistant", Content: textReply,
			},
		)
	}
	return textReply
}

// Talk 取得带 CQ 码的回复消息
func (l *LolimiAi2) Talk(_ int64, msg, nickname string) string {
	return l.TalkPlain(0, msg, nickname)
}

// Talk 取得带 CQ 码的回复消息
func (l *LolimiAi2Mem) Talk(_ int64, msg, nickname string) string {
	return l.TalkPlain(0, msg, nickname)
}
