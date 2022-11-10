package aireply

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
)

// XiaoAiReply 小爱回复类
type XiaoAiReply struct{}

const (
	xiaoaiURL     = "http://81.70.100.130/api/xiaoai.php?n=text&msg=%v"
	xiaoaiBotName = "小爱"
)

func (*XiaoAiReply) String() string {
	return "小爱"
}

// TalkPlain 取得回复消息
func (*XiaoAiReply) TalkPlain(msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, xiaoaiBotName)
	u := fmt.Sprintf(xiaoaiURL, url.QueryEscape(msg))
	replyMsg, err := web.GetData(u)
	if err != nil {
		return "ERROR:" + err.Error()
	}
	textReply := strings.ReplaceAll(binary.BytesToString(replyMsg), xiaoaiBotName, nickname)
	if textReply == "" {
		textReply = nickname + "听不懂你的话了, 能再说一遍吗"
	}
	textReply = strings.ReplaceAll(textReply, "小米智能助理", "电子宠物")

	return textReply
}

// Talk 取得带 CQ 码的回复消息
func (x *XiaoAiReply) Talk(msg, nickname string) string {
	return x.TalkPlain(msg, nickname)
}
