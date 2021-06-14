package pixiv

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

// P站 无污染 IP 地址
var IPTables = map[string]string{
	"pixiv.net":   "210.140.131.223:443",
	"i.pximg.net": "210.140.92.142:443",
}

// Works 获取插画信息
func Works(id int) ([]byte, error) {
	return netPost(fmt.Sprintf("https://pixiv.net/ajax/illust/%d", id))
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
