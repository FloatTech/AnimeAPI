// Package aireply 人工智能回复
package aireply

import "fmt"

// AIReply 公用智能回复类
type AIReply interface {
	// Talk 取得带 CQ 码的回复消息
	Talk(uid int64, msg, nickname string) string
	// Talk 取得文本回复消息
	TalkPlain(uid int64, msg, nickname string) string
	// String 获得实际使用的回复服务名
	fmt.Stringer
}
