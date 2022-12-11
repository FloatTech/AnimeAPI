package chatgpt

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	"github.com/google/uuid"
)

const (
	SESSION_TOKEN = "__Secure-next-auth.session-token"
	UA            = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"
)

var (
	ErrRequestTooFast  = errors.New("request too fast")
	ErrEmptyResponse   = errors.New("empty response")
	ErrNilSessionToken = errors.New("refresh session failed: nil token")
	ErrNilAuth         = errors.New("refresh session failed: nil auth")
)

type ChatGPT struct {
	config *Config
	Auth   string
	ConvID string
	ParnID string
}

func NewChatGPT(config *Config) *ChatGPT {
	return &ChatGPT{config: config}
}

func (c *ChatGPT) id() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func (c *ChatGPT) setchatheaders(req *http.Request) {
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+c.Auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://chat.openai.com")
	req.Header.Set("Referer", "https://chat.openai.com/chat")
}

func (c *ChatGPT) getbody(prompts ...string) *bytes.Buffer {
	body := bytes.NewBuffer(make([]byte, 0, 4096))
	body.WriteString(`{"action":"next","messages":[{"id":"`)
	body.WriteString(c.id())
	body.WriteString(`","role":"user","content":{"content_type":"text","parts":`)
	_ = json.NewEncoder(body).Encode(&prompts)
	body.Truncate(body.Len() - 1)
	switch {
	case c.ConvID != "":
		body.WriteString(`}}],"conversation_id":"`)
		body.WriteString(c.ConvID)
		body.WriteString(`","parent_message_id":"`)
		body.WriteString(c.ParnID)
		body.WriteByte('"')
	case c.ParnID != "":
		body.WriteString(`}}],"parent_message_id":"`)
		body.WriteString(c.ParnID)
		body.WriteByte('"')
	default:
		body.WriteString(`}}]`)
	}
	body.WriteString(`,"model":"text-davinci-002-render"}`)
	return body
}

type chatresponse struct {
	Message struct {
		ID      string `json:"id"`
		Content struct {
			ContentType string   `json:"content_type"`
			Parts       []string `json:"parts"`
		} `json:"content"`
		Weight float64 `json:"weight"`
	} `json:"message"`
	ConversationID string `json:"conversation_id"`
	Error          any    `json:"error"`
}

func (c *ChatGPT) GetChatResponse(prompt string) (string, error) {
	if c.Auth == "" {
		err := c.RefreshSession()
		if err != nil {
			return "", err
		}
	}
	if c.ParnID == "" {
		c.ParnID = uuid.Must(uuid.NewRandom()).String()
	}
	body := c.getbody(prompt)
	req, err := http.NewRequest("POST", API+"backend-api/conversation", body)
	if err != nil {
		return "", err
	}
	c.setchatheaders(req)
	req.Header.Set("Content-Length", strconv.Itoa(body.Len()))
	req.Header.Set("User-Agent", c.config.UA)
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS12,
			},
		},
		Timeout: c.config.Timeout,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		return "", ErrRequestTooFast
	}
	s := bufio.NewScanner(resp.Body)
	lastline := ""
	line := ""
	for s.Scan() {
		l := s.Text()
		if l == "" {
			continue
		}
		lastline = line
		line = l
	}
	if len(lastline) <= 6 {
		return "", ErrEmptyResponse
	}
	var rsp chatresponse
	err = json.Unmarshal(binary.StringToBytes(lastline[6:]), &rsp)
	if err != nil {
		return "", err
	}
	if rsp.Error != nil {
		return "", errors.New(fmt.Sprint(rsp.Error))
	}
	if len(rsp.Message.Content.Parts) == 0 || rsp.ConversationID == "" || rsp.Message.ID == "" {
		return "", ErrEmptyResponse
	}
	c.ConvID = rsp.ConversationID
	c.ParnID = rsp.Message.ID
	return rsp.Message.Content.Parts[0], nil
}

func (c *ChatGPT) RefreshSession() error {
	req, err := http.NewRequest("GET", API+"api/auth/session", nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{Name: SESSION_TOKEN, Value: c.config.SessionToken})
	req.Header.Set("User-Agent", c.config.UA)
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS12,
			},
		},
		Timeout: c.config.Timeout,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var rsp struct {
		Token string `json:"accessToken"`
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == SESSION_TOKEN {
			if cookie.Value == "" {
				return ErrNilSessionToken
			}
			c.config.SessionToken = cookie.Value
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil {
		return err
	}
	if rsp.Token == "" {
		return ErrNilAuth
	}
	c.Auth = rsp.Token
	return nil
}
