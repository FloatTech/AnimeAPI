// Package imgpool 图片缓存池
package imgpool

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FloatTech/zbputils/pool"
	"github.com/FloatTech/zbputils/process"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const cacheurl = "https://gchat.qpic.cn/gchatpic_new//%s/0"

type Image struct {
	img  *pool.Item
	n, f string
}

// GetImage name
func GetImage(name string) (m *Image, err error) {
	m = new(Image)
	m.n = name
	m.img, err = pool.GetItem(name)
	if err == nil && m.img.String() != "" {
		_, err = http.Head(m.String())
		if err != nil {
			err = errors.New("img file outdated")
			return
		}
		return
	}
	err = errors.New("no such img")
	return
}

// NewImage context name file
func NewImage(ctx *zero.Ctx, name, f string) (m *Image, err error) {
	m = new(Image)
	m.n = name
	m.SetFile(f)
	m.img, err = pool.GetItem(name)
	if err == nil && m.img.String() != "" {
		_, err = http.Head(m.String())
		if err == nil {
			return
		}
	}
	err = m.Push(ctx)
	return
}

// String url
func (m *Image) String() string {
	return fmt.Sprintf(cacheurl, m.img.String())
}

// SetFile f
func (m *Image) SetFile(f string) {
	if strings.HasPrefix(f, "http") {
		m.f = f
	} else {
		m.f = "file:///" + f
	}
}

func (m *Image) Push(ctx *zero.Ctx) (err error) {
	id := ctx.SendPrivateMessage(ctx.Event.SelfID, message.Message{message.Image(m.f)})
	if id == 0 {
		err = errors.New("send image error")
		return
	}
	defer process.SleepAbout1sTo2s() // 防止风控
	msg := ctx.GetMessage(message.NewMessageID(strconv.FormatInt(id, 10)))
	for _, e := range msg.Elements {
		if e.Type == "image" {
			u := e.Data["url"]
			u = u[:strings.LastIndex(u, "/")]
			u = u[strings.LastIndex(u, "/")+1:]
			if u != "" {
				m.img, err = pool.NewItem(m.n, u)
				logrus.Infoln("[imgpool] 缓存:", m.n, "url:", u)
				_ = m.img.Push("minamoto")
			} else {
				err = errors.New("get msg error")
			}
			return
		}
	}
	err = errors.New("get msg error")
	return
}
