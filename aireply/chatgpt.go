package aireply

import (
	"strings"
	"sync"
	"time"

	"github.com/FloatTech/AnimeAPI/chatgpt"
	"github.com/FloatTech/ttl"
	"github.com/RomiChan/syncx"
	"github.com/sirupsen/logrus"
)

// ChatGPT 回复类
type ChatGPT struct {
	sync.Mutex
	c *chatgpt.Config
	s *ttl.Cache[int64, *chatgpt.ChatGPT]
	m *syncx.Map[int64, struct{}]
}

func NewChatGPT(config *chatgpt.Config) (c *ChatGPT) {
	c = &ChatGPT{
		c: config,
		m: &syncx.Map[int64, struct{}]{},
	}
	c.s = ttl.NewCacheOn(time.Hour, [4]func(int64, *chatgpt.ChatGPT){
		func(uid int64, chat *chatgpt.ChatGPT) {
			go func() {
				for range time.NewTicker(c.c.RefreshInterval).C {
					if _, ok := c.m.Load(uid); ok {
						c.m.Delete(uid)
						return
					}
					err := chat.RefreshSession()
					if err != nil {
						logrus.Errorln("[chatgpt] 刷新 session 错误:", err)
					}
				}
			}()
		}, nil,
		func(uid int64, _ *chatgpt.ChatGPT) {
			c.m.Store(uid, struct{}{})
		}, nil,
	})
	return
}

func (c *ChatGPT) String() string {
	return "ChatGPT"
}

func (c *ChatGPT) Talk(uid int64, msg, nickname string) string {
	c.Lock()
	defer c.Unlock()
	chat := c.s.Get(uid)
	if chat == nil {
		chat = chatgpt.NewChatGPT(c.c)
		c.s.Set(uid, chat)
	}
	msg = strings.ReplaceAll(msg, nickname, "你")
	reply, err := chat.GetChatResponse(msg)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return reply
}

func (c *ChatGPT) TalkPlain(uid int64, msg, nickname string) string {
	return c.Talk(uid, msg, nickname)
}

func (c *ChatGPT) Reset(uid int64) {
	c.s.Delete(uid)
}
