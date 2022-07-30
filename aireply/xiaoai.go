package aireply

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
)

// XiaoAiReply 小爱回复类
type XiaoAiReply struct{}

const (
	xiaoaiURL     = "https://yang520.ltd/api/xiaoai.php?msg=%v"
	xiaoaiBotName = "小爱"
)

func (*XiaoAiReply) String() string {
	return "小爱"
}

// TalkPlain 取得回复消息
func (*XiaoAiReply) TalkPlain(msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, xiaoaiBotName)
	u := fmt.Sprintf(xiaoaiURL, url.QueryEscape(msg))
	data, err := web.RequestDataWith(web.NewDefaultClient(), u, "GET", "", web.RandUA())
	if err != nil {
		return "ERROR:" + err.Error()
	}
	replystr := gjson.Get(binary.BytesToString(data), "text").String()
	textReply := strings.ReplaceAll(replystr, xiaoaiBotName, nickname)
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
