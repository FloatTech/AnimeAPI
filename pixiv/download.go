package pixiv

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Download(link, filedir, filename string) error {
	// 取文件路径
	if strings.Contains(filedir, `/`) && !strings.HasSuffix(filedir, `/`) {
		filedir += `/`
	}
	if strings.Contains(filedir, `\`) && !strings.HasSuffix(filedir, `\`) {
		filedir += `\`
	}
	filepath := filedir + filename
	// 路径目录不存在则创建目录
	if _, err := os.Stat(filedir); err != nil && !os.IsExist(err) {
		if err := os.MkdirAll(filedir, 0644); err != nil {
			return err
		}
	}
	// P站特殊客户端
	client := &http.Client{
		// 解决中国大陆无法访问的问题
		Transport: &http.Transport{
			DisableKeepAlives: false,
			// 隐藏 sni 标志
			TLSClientConfig: &tls.Config{
				ServerName:         "-",
				InsecureSkipVerify: true,
			},
			// 更改 dns
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp", IPTables["i.pximg.net"])
			},
		},
	}
	// 网络请求
	request, _ := http.NewRequest("GET", link, nil)
	request.Header.Set("Host", "i.pximg.net")
	request.Header.Set("Referer", "https://www.pixiv.net/")
	request.Header.Set("Accept", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 验证接收到的长度
	length, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	data, _ := ioutil.ReadAll(resp.Body)
	if length != len(data) {
		return errors.New("download not complete")
	}
	// 获取文件后缀
	switch resp.Header.Get("Content-Type") {
	case "image/jpeg":
		filepath += ".jpg"
	case "image/png":
		filepath += ".png"
	case "image/gif":
		filepath += ".gif"
	default:
		filepath += ".jpg"
	}
	// 写入文件
	f, _ := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)
	return nil
}
