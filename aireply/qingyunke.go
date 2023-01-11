package aireply

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

// QYK 青云客回复类
type QYK struct {
	u string
	n string
	b []string
}

const (
	QYKURL     = "http://api.qingyunke.com/api.php?key=free&appid=0&msg=%v"
	QYKBotName = "菲菲"
)

var (
	qykMatchFace = regexp.MustCompile(`\{face:(\d+)\}(.*)`)
)

func NewQYK(u, name string, banwords ...string) *QYK {
	return &QYK{u: u, n: name, b: banwords}
}

func (*QYK) String() string {
	return "青云客"
}

// Talk 取得带 CQ 码的回复消息
func (q *QYK) Talk(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, q.n)
	u := fmt.Sprintf(q.u, url.QueryEscape(msg))
	data, err := web.RequestDataWith(web.NewDefaultClient(), u, "GET", "", web.RandUA(), nil)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	replystr := gjson.Get(binary.BytesToString(data), "content").String()
	replystr = strings.ReplaceAll(replystr, "{face:", "[CQ:face,id=")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, "}", "]")
	replystr = strings.ReplaceAll(replystr, q.n, nickname)
	for _, w := range q.b {
		if strings.Contains(replystr, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	return replystr
}

// TalkPlain 取得回复消息
func (q *QYK) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, q.n)

	u := fmt.Sprintf(q.u, url.QueryEscape(msg))
	data, err := web.RequestDataWith(web.NewDefaultClient(), u, "GET", "", web.RandUA(), nil)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	replystr := gjson.Get(binary.BytesToString(data), "content").String()
	replystr = qykMatchFace.ReplaceAllLiteralString(replystr, "")
	replystr = strings.ReplaceAll(replystr, "{br}", "\n")
	replystr = strings.ReplaceAll(replystr, q.n, nickname)
	for _, w := range q.b {
		if strings.Contains(replystr, w) {
			return "ERROR: 回复可能含有敏感内容"
		}
	}
	return replystr
}
