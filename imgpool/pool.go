// Package imgpool 图片缓存池
package imgpool

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/pool"
	"github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const cacheurl = "https://gchat.qpic.cn/gchatpic_new//%s/0"

var (
	ErrImgFileOutdated = errors.New("img file outdated")
	ErrNoSuchImg       = errors.New("no such img")
	ErrSendImg         = errors.New("send image error")
	ErrGetMsg          = errors.New("get msg error")
)

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
			err = ErrImgFileOutdated
			return
		}
		return
	}
	err = ErrNoSuchImg
	return
}

// NewImage context name file
func NewImage(send ctxext.NoCtxSendMsg, get ctxext.NoCtxGetMsg, name, f string) (m *Image, hassent bool, err error) {
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
	hassent, err = m.Push(send, get)
	return
}

// String url
func (m *Image) String() string {
	if m.img == nil {
		return m.f
	}
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

func (m *Image) Push(send ctxext.NoCtxSendMsg, get ctxext.NoCtxGetMsg) (hassent bool, err error) {
	id := send(message.Message{message.Image(m.f)})
	if id == 0 {
		err = ErrSendImg
		return
	}
	hassent = true
	msg := get(id)
	for _, e := range msg.Elements {
		if e.Type == "image" {
			u := e.Data["url"]
			i := strings.LastIndex(u, "/")
			if i <= 0 {
				break
			}
			u = u[:i]
			i = strings.LastIndex(u, "/")
			if i <= 0 {
				break
			}
			u = u[i+1:]
			if u == "" {
				break
			}
			m.img, err = pool.NewItem(m.n, u)
			logrus.Infoln("[imgpool] 缓存:", m.n, "url:", u)
			_ = m.img.Push("minamoto")
			return
		}
	}
	err = ErrGetMsg
	return
}
