// Package pixiv (copied from plugin/pixiv/api/client.go)
package pixiv

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/FloatTech/AnimeAPI/pixiv/model"
)

// HTTPStatusError ...
type HTTPStatusError struct {
	StatusCode int
	URL        string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("下载图片失败: HTTP %d", e.StatusCode)
}

// Client 封装 HTTP 客户端与 Pixiv 请求逻辑
type Client struct {
	*http.Client
	transport *http.Transport
}

// NewClient ...
func NewClient() *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{MaxVersion: tls.VersionTLS13},
		Proxy:           http.ProxyFromEnvironment,
	}
	return &Client{
		Client: &http.Client{
			Transport: transport,
			Timeout:   time.Minute,
		},
		transport: transport,
	}
}

// SetProxy 设置代理
func (c *Client) SetProxy(proxyURL string) error {
	if c == nil || c.transport == nil {
		return errors.New("pixiv client is nil")
	}
	proxyURL = strings.TrimSpace(proxyURL)
	if proxyURL == "" {
		c.transport.Proxy = http.ProxyFromEnvironment
		return nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	c.transport.Proxy = http.ProxyURL(u)
	return nil
}

// SearchPixivIllustrations ...
func (c *Client) SearchPixivIllustrations(accessToken, url string) (*model.RootEntity, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 (Android 11; Pixel 5)")
	req.Header.Set("App-OS", "android")
	req.Header.Set("App-OS-Version", "11")
	req.Header.Set("App-Version", "5.0.234")

	req.Header.Set("Accept-Language", "en_US")
	req.Header.Set("Referer", "https://app-api.pixiv.net/")
	req.Header.Set("Connection", "keep-alive")

	req.Host = "app-api.pixiv.net"

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New("搜索失败: " + resp.Status + "\nbody: " + string(body))
	}

	var result model.RootEntity
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) fetchOnce(targetURL, referer string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 (Android 11; Pixel 5)")

	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, errors.New("请求失败: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return data, resp.StatusCode, nil
}

// FetchPixivImage 直接从 Pixiv 下载图片
func (c *Client) FetchPixivImage(illust model.IllustCache, url string) ([]byte, error) {
	fmt.Println("下载", illust.PID)

	if c == nil {
		fmt.Println("FetchPixivImage called on nil IllustCache")
		return nil, nil
	}

	data, status, err := c.fetchOnce(url, "https://www.pixiv.net/")
	if err == nil && status == http.StatusOK {
		return data, nil
	}
	if status == http.StatusNotFound {
		return nil, &HTTPStatusError{StatusCode: status, URL: url}
	}

	if err != nil {
		return nil, err
	}

	if status != 0 {
		return nil, &HTTPStatusError{StatusCode: status, URL: url}
	}

	return nil, errors.New("下载图片失败")
}
