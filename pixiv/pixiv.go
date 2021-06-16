package pixiv

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

// P站 无污染 IP 地址
var IPTables = map[string]string{
	"pixiv.net":   "210.140.131.223:443",
	"i.pximg.net": "210.140.92.142:443",
}

//插画结构体
type ImgAll struct {
	ID          string
	Title       string
	Type        string
	Description string
	Small       string
	Original    string
	AuthorID    string
	AuthorName  string
	Width       int64
	Height      int64
}

//[]byte转tjson类型
type tjson []byte

//解析json
func (data tjson) Get(path string) gjson.Result {
	return gjson.Get(string(data), path)
}

// Works 获取插画信息
func Works(id string) (i ImgAll, err error) {
	var b []byte
	b, err = netPost(fmt.Sprintf("https://pixiv.net/ajax/illust/%s", id))
	if err != nil {
		return
	}
	body := tjson(b)
	i.ID = id
	i.Title = body.Get("body.illustTitle").String()
	i.Type = body.Get("body.illustType").String()
	i.Description = body.Get("body.description").String()
	i.Small = body.Get("body.urls.small").String()
	i.Original = body.Get("body.urls.original").String()
	i.AuthorID = body.Get("body.userId").String()
	i.AuthorName = body.Get("body.userName").String()
	i.Width = body.Get("body.width").Int()
	i.Height = body.Get("body.height").Int()
	return
}

// Download 下载图片
func Download(link string) ([]byte, error) {
	return netPost(link)
}

// netPost 返回请求数据
func netPost(link string) ([]byte, error) {
	// 获取IP地址
	domain, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	IP := IPTables[domain.Host]
	// P站特殊客户端
	client := &http.Client{
		// 解决中国大陆无法访问的问题
		Transport: &http.Transport{
			DisableKeepAlives: true,
			// 隐藏 sni 标志
			TLSClientConfig: &tls.Config{
				ServerName:         "-",
				InsecureSkipVerify: true,
			},
			// 更改 dns
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp", IP)
			},
		},
	}
	// 网络请求
	request, _ := http.NewRequest("POST", link, nil)
	request.Header.Set("Host", domain.Host)
	request.Header.Set("Referer", "https://www.pixiv.net/")
	request.Header.Set("Accept", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	return result, nil
}
