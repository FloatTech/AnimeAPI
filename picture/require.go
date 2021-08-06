package picture

import (
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// CmdMatch 命令匹配
func CmdMatch() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		for _, elem := range ctx.Event.Message {
			if elem.Type == "text" {
				text := strings.ReplaceAll(elem.Data["text"], " ", "")
				if text != ctx.State["keyword"].(string) {
					return false
				}
			}
		}
		return true
	}
}

// Exists 消息含有图片返回 true
func Exists() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		var urls = []string{}
		for _, elem := range ctx.Event.Message {
			if elem.Type == "image" {
				urls = append(urls, elem.Data["url"])
			}
		}
		if len(urls) > 0 {
			ctx.State["image_url"] = urls
			return true
		}
		return false
	}
}

// MustGiven 消息不存在图片阻塞60秒至有图片，超时返回 false
func MustGiven() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		if Exists()(ctx) {
			return true
		}
		// 没有图片就索取
		ctx.SendChain(message.Text("请发送一张图片"))
		next := zero.NewFutureEvent("message", 999, false, zero.CheckUser(ctx.Event.UserID), Exists())
		recv, cancel := next.Repeat()
		select {
		case <-time.After(time.Second * 120):
			return false
		case e := <-recv:
			cancel()
			newCtx := &zero.Ctx{Event: e, State: zero.State{}}
			if Exists()(newCtx) {
				ctx.State["image_url"] = newCtx.State["image_url"]
			}
			return true
		}
	}
}
