package bilibili

import (
	"regexp"
	"strconv"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var re = regexp.MustCompile(`^\d+$`)

func RequireUser(cfg *CookieConfig) func(ctx *zero.Ctx) bool {
	return func(ctx *zero.Ctx) bool {
		keyword := ctx.State["regex_matched"].([]string)[1]
		if !re.MatchString(keyword) {
			searchRes, err := cfg.SearchUser(keyword)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			ctx.State["uid"] = strconv.FormatInt(searchRes[0].Mid, 10)
			return true
		}
		next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession())
		recv, cancel := next.Repeat()
		defer cancel()
		ctx.SendChain(message.Text("输入为纯数字, 请选择查询uid还是用户名, 输入对应序号：\n0. 查询uid\n1. 查询用户名"))
		for {
			select {
			case <-time.After(time.Second * 10):
				ctx.SendChain(message.Text("时间太久啦！", zero.BotConfig.NickName[0], "帮你选择查询uid"))
				ctx.State["uid"] = keyword
				return true
			case c := <-recv:
				msg := c.Event.Message.ExtractPlainText()
				num, err := strconv.Atoi(msg)
				if err != nil {
					ctx.SendChain(message.Text("请输入数字!"))
					continue
				}
				if num < 0 || num > 1 {
					ctx.SendChain(message.Text("序号非法!"))
					continue
				}
				if num == 0 {
					ctx.State["uid"] = keyword
					return true
				} else if num == 1 {
					searchRes, err := cfg.SearchUser(keyword)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return false
					}
					ctx.State["uid"] = strconv.FormatInt(searchRes[0].Mid, 10)
					return true
				}
			}
		}
	}
}
