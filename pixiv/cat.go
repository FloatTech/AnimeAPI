package pixiv

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/web"
)

// Generate API 返回结果
type Generate struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error"`
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Artist   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"artist"`
	Multiple          bool     `json:"multiple"`
	OriginalURL       string   `json:"original_url"`
	OriginalURLProxy  string   `json:"original_url_proxy"`
	OriginalUrls      []string `json:"original_urls"`
	OriginalUrlsProxy []string `json:"original_urls_proxy"`
	Thumbnails        []string `json:"thumbnails"`
}

// Cat 调用 pixiv.cat 的 generate API
func Cat(id int64) (*Generate, error) {
	data, err := web.RequestDataWithHeaders(&http.Client{
		Transport: &http.Transport{
			DialTLS: func(_, _ string) (net.Conn, error) {
				return tls.Dial("tcp", "66.42.35.2:443", &tls.Config{
					ServerName: "api.pixiv.cat",
					MaxVersion: tls.VersionTLS12,
				})
			},
		}},
		"https://api.pixiv.cat/v1/generate",
		"POST",
		func(r *http.Request) error {
			r.Header.Set("authority", "api.pixiv.cat")
			r.Header.Set("accept", "*/*")
			r.Header.Set("accept-language", "zh,zh-CN;q=0.9,zh-HK;q=0.8,zh-TW;q=0.7,ja;q=0.6,en;q=0.5,en-GB;q=0.4,en-US;q=0.3")
			r.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
			r.Header.Set("origin", "https://pixiv.cat")
			r.Header.Set("referer", "https://pixiv.cat/")
			r.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.76")
			return nil
		},
		strings.NewReader("p="+url.QueryEscape("https://www.pixiv.net/member_illust.php?illust_id="+strconv.FormatInt(id, 10))),
	)
	if err != nil {
		return nil, err
	}
	g := &Generate{}
	err = json.Unmarshal(data, g)
	if err != nil {
		return nil, err
	}
	return g, nil
}
