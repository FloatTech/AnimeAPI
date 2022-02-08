package aireply

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
)

// QYKReply 青云客回复类
type QYKReply struct{}

const (
	qykURL     = "http://api.qingyunke.com/api.php?key=free&appid=0&msg=%s"
	qykBotName = "菲菲"
)

var (
	qykMatchFace = regexp.MustCompile(`\{face:(\d+)\}(.*)`)
)

func (*QYKReply) String() string {
	return "青云客"
}

// Talk 取得带 CQ 码的回复消息
func (*QYKReply) Talk(msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, qykBotName)

	u := fmt.Sprintf(qykURL, url.QueryEscape(msg))
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "ERROR:" + err.Error()
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.qingyunke.com")
	resp, err := client.Do(req)
	if err != nil {
		return "ERROR:" + err.Error()
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "ERROR:" + err.Error()
	}

	replystr := gjson.Get(binary.BytesToString(bytes), "content").String()
	replystr = strings.ReplaceAll(replystr, "{face:", "[CQ:face,id=")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, "}", "]")
	replystr = strings.ReplaceAll(replystr, qykBotName, nickname)

	return replystr
}

// TalkPlain 取得回复消息
func (*QYKReply) TalkPlain(msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, qykBotName)

	u := fmt.Sprintf(qykURL, url.QueryEscape(msg))
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	// 自定义Header
	req.Header.Set("User-Agent", web.RandUA())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api.qingyunke.com")
	resp, err := client.Do(req)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "ERROR: " + err.Error()
	}

	replystr := gjson.Get(binary.BytesToString(bytes), "content").String()
	replystr = qykMatchFace.ReplaceAllLiteralString(replystr, "")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, qykBotName, nickname)

	return replystr
}
