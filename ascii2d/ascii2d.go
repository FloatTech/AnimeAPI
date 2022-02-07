package ascii2d

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	xpath "github.com/antchfx/htmlquery"
)

type Result struct {
	Info   string // Info 图片分辨率 格式 大小信息
	Link   string // Link 图片链接
	Name   string // Name 图片名
	Author string // Author 作者链接
	AuthNm string // AuthNm 作者名
	Thumb  string // Thumb 缩略图链接
	Type   string // Type pixiv / twitter ...
}

func Ascii2d(image string) (r []*Result, err error) {
	var (
		api = "https://ascii2d.net/search/uri"
	)
	transport := http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS12,
		},
	}
	client := &http.Client{
		Transport: &transport,
	}

	// 包装请求参数
	data := url.Values{}
	data.Set("uri", image) // 图片链接
	fromData := strings.NewReader(data.Encode())

	// 网络请求
	req, _ := http.NewRequest("POST", api, fromData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析XPATH
	doc, err := xpath.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	// 取出每个返回的结果
	list := xpath.Find(doc, `//div[@class="row item-box"]`)
	if len(list) == 0 {
		return
	}
	r = make([]*Result, 0, len(list))
	// 遍历结果
	for _, n := range list {
		linkPath := xpath.FindOne(n, `//div[2]/div[3]/h6/a[1]`)
		authPath := xpath.FindOne(n, `//div[2]/div[3]/h6/a[2]`)
		picPath := xpath.FindOne(n, `//div[1]/img`)
		if linkPath != nil && authPath != nil && picPath != nil {
			r = append(r, &Result{
				Info:   xpath.InnerText(xpath.FindOne(n, `//div[2]/small`)),
				Link:   xpath.SelectAttr(linkPath, "href"),
				Name:   xpath.InnerText(linkPath),
				Author: xpath.SelectAttr(authPath, "href"),
				AuthNm: xpath.InnerText(authPath),
				Thumb:  "https://ascii2d.net" + xpath.SelectAttr(picPath, "src"),
				Type:   strings.Trim(xpath.InnerText(xpath.FindOne(n, `//div[2]/div[3]/h6/small`)), "\n"),
			})
		}
	}
	return
}
