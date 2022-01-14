package aireply

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// XiaoAiReply 小爱回复类
type XiaoAiReply struct{}

const (
	xiaoaiURL     = "http://81.70.100.130/api/xiaoai.php?msg=%s&n=text"
	xiaoaiBotName = "小爱"
)

func (*XiaoAiReply) String() string {
	return "小爱"
}

// TalkPlain 取得回复消息
func (*XiaoAiReply) TalkPlain(msg string) string {
	msg = strings.ReplaceAll(msg, zero.BotConfig.NickName[0], xiaoaiBotName)

	u := fmt.Sprintf(xiaoaiURL, url.QueryEscape(msg))
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "81.70.100.130")
	resp, err := client.Do(req)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	replystr := helper.BytesToString(bytes)
	textReply := strings.ReplaceAll(replystr, xiaoaiBotName, zero.BotConfig.NickName[0])
	if textReply == "" {
		textReply = zero.BotConfig.NickName[0] + "听不懂你的话了，能再说一遍吗"
	}
	textReply = strings.ReplaceAll(textReply, "小米智能助理", "电子宠物")

	return textReply
}

// Talk 取得回复消息
func (x *XiaoAiReply) Talk(msg string) message.Message {
	return message.Message{message.Text(x.TalkPlain(msg))}
}
