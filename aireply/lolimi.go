package aireply

import (
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
	k string
	n string
	b []string
}

const (
	lolimiURL = "https://apii.lolimi.cn"
	// MomoURL api地址
	MomoURL = lolimiURL + "/api/mmai/mm?key=%v&msg=%v"
	// MomoBotName ...
	MomoBotName = "沫沫"
	// JingfengURL api地址
	JingfengURL = lolimiURL + "/api/jjai/jj?key=%v&msg=%v"
	// JingfengBotName ...
	JingfengBotName = "婧枫"
)

// NewLolimiAi ...
func NewLolimiAi(u, k, name string, banwords ...string) *LolimiAi {
	return &LolimiAi{u: u, k: k, n: name, b: banwords}
}

// String ...
func (l *LolimiAi) String() string {
	return l.n
}

// TalkPlain 取得回复消息
func (l *LolimiAi) TalkPlain(_ int64, msg, nickname string) string {
	msg = strings.ReplaceAll(msg, nickname, l.n)
	u := fmt.Sprintf(l.u, l.k, url.QueryEscape(msg))
	data, err := web.GetData(u)
	if err != nil {
		return "ERROR: " + strings.ReplaceAll(err.Error(), l.k, "********")
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
