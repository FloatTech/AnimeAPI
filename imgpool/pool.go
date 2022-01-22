// Package imgpool 图片缓存池
package imgpool

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/zbputils/pool"
	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const cacheurl = "https://gchat.qpic.cn/gchatpic_new//%s/0"
const imgpoolgrp = 117500479

type Image struct {
	img *pool.Item
}

// NewImage context name file
func NewImage(ctx *zero.Ctx, name, f string) (m Image, err error) {
	var data []byte
	if strings.HasPrefix(f, "http") {
		data, err = web.GetData(f)
	} else {
		data, err = os.ReadFile(f)
	}
	if err != nil {
		return
	}
	m.img, err = pool.GetItem(name)
	if err == nil {
		return
	}
	id := ctx.SendGroupMessage(imgpoolgrp, message.Message{message.Text(name), message.Image("base64://" + base64.StdEncoding.EncodeToString(data))})
	if id == 0 {
		err = errors.New("send image error")
		return
	}
	msg := ctx.GetMessage(message.NewMessageID(strconv.FormatInt(id, 10)))
	for _, e := range msg.Elements {
		if e.Type == "image" {
			u := e.Data["file"]
			u = u[:strings.LastIndex(u, "/")]
			u = u[strings.LastIndex(u, "/")+1:]
			m.img, err = pool.NewItem(name, u)
			break
		}
	}
	return
}

// RegisterListener key engine
func RegisterListener(key string, en *zero.Engine) {
	en.OnMessage(zero.OnlyGroup, func(ctx *zero.Ctx) bool {
		return ctx.Event.GroupID == imgpoolgrp && ctx.Event.MessageType == "image"
	}).SetBlock(true).FirstPriority().
		Handle(func(ctx *zero.Ctx) {
			var u, n string
			for _, e := range ctx.Event.Message {
				if e.Type == "image" {
					u = e.Data["file"]
					u = u[:strings.LastIndex(u, "/")]
					u = u[strings.LastIndex(u, "/")+1:]
				} else if e.Type == "text" {
					n = e.Data["text"]
				}
			}
			img, err := pool.NewItem(n, u)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			err = img.Push(key)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
		})
}

// String url
func (m Image) String() string {
	return fmt.Sprintf(cacheurl, m.img.String())
}
