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
	u string          // API 地址
	n string          // AI 名称
	k string          // API 密钥
	b []string        // Banwords
	t bool            // API 响应模式是否为文本
	l int             // 记忆限制数（小于 1 的值可以禁用记忆模式）
	m []lolimiMessage // 记忆数据
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
func NewLolimiAi(u, name string, key string, textMode bool, memoryLimit int, banwords ...string) *LolimiAi {
	return &LolimiAi{u: u, n: name, k: key, t: textMode, l: memoryLimit, b: banwords}
}

// String ...
func (l *LolimiAi) String() string {
	return l.n
}

// TalkPlain 取得回复消息
func (l *LolimiAi) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, l.n)
	var u string
	var data []byte
	var err error
	if l.l > 0 {
		u = fmt.Sprintf(l.u, url.QueryEscape(l.k))
		json, err := json.Marshal(
			append(l.m,
				lolimiMessage{
					Role: "user", Content: msg,
				},
			),
		)
		if err != nil {
			return "ERROR: " + err.Error()
		}
		// TODO: 可能会返回
		// "请使用psot格式请求如有疑问进官方群"
		data, err = web.PostData(u, "application/json", bytes.NewReader(json))
		if err != nil {
			return "ERROR: " + err.Error()
		}
	} else {
		u := fmt.Sprintf(l.u, url.QueryEscape(l.k), url.QueryEscape(msg))
		data, err = web.GetData(u)
	}
	if err != nil {
		errMsg := err.Error()
		// Remove the key from error message
		errMsg = strings.ReplaceAll(errMsg, l.k, "********")
		return "ERROR: " + errMsg
	}
	var replystr string
	if l.t {
		replystr = binary.BytesToString(data)
	} else {
		replystr = gjson.Get(binary.BytesToString(data), "data.output").String()
	}
	// TODO: 是否要删除遗留代码
	replystr = strings.ReplaceAll(replystr, "<img src=\"", "[CQ:image,file=")
	replystr = strings.ReplaceAll(replystr, "<br>", "\n")
	replystr = strings.ReplaceAll(replystr, "\" />", "]")
	textReply := strings.ReplaceAll(replystr, l.n, nickname)
	for _, w := range l.b {
		if strings.Contains(textReply, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	if l.l > 0 {
		// 添加记忆
		if len(l.m) >= l.l-1 && len(l.m) >= 2 {
			l.m = l.m[2:]
		}
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
func (l *LolimiAi) Talk(_ int64, msg, nickname string) string {
	return l.TalkPlain(0, msg, nickname)
}
