package aireply

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
)

// XiaoAi 小爱回复类
type XiaoAi struct {
	u string
	n string
	b []string
}

const (
	XiaoAiURL     = "http://81.70.100.130/api/xiaoai.php?n=text&msg=%v"
	XiaoAiBotName = "小爱"
)

func NewXiaoAi(u, name string, banwords ...string) *XiaoAi {
	return &XiaoAi{u: u, n: name, b: banwords}
}

func (*XiaoAi) String() string {
	return "小爱"
}

// TalkPlain 取得回复消息
func (x *XiaoAi) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, x.n)
	u := fmt.Sprintf(x.u, url.QueryEscape(msg))
	replyMsg, err := web.GetData(u)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	textReply := strings.ReplaceAll(binary.BytesToString(replyMsg), x.n, nickname)
	if textReply == "" {
		textReply = nickname + "听不懂你的话了, 能再说一遍吗"
	}
	textReply = strings.ReplaceAll(textReply, "小米智能助理", "电子宠物")
	for _, w := range x.b {
		if strings.Contains(textReply, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	return textReply
}

// Talk 取得带 CQ 码的回复消息
func (x *XiaoAi) Talk(_ int64, msg, nickname string) string {
	return x.TalkPlain(0, msg, nickname)
}
