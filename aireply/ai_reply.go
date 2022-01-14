// Package aireply 人工智能回复
package aireply

import (
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	modeMap = func() (m map[string]AIReply) {
		setReplyMap := func(m map[string]AIReply, r AIReply) {
			m[r.String()] = r
		}
		m = make(map[string]AIReply, 2)
		setReplyMap(m, &QYKReply{})
		setReplyMap(m, &XiaoAiReply{})
		return
	}()
)

// AIReply 公用智能回复类
type AIReply interface {
	// Talk 取得回复消息
	Talk(string) message.Message
	// Talk 取得文本回复消息
	TalkPlain(string) string
	// String 获得模式
	String() string
}

// NewAIReply 智能回复简单工厂
func NewAIReply(mode string) AIReply {
	return modeMap[mode]
}
