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

var (
	// P站特殊客户端
	client = &http.Client{
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
)

// DownloadData 下载 link，返回 &data, suffix, error
func DownloadData(link string) (*[]byte, string, error) {
	// 网络请求
	request, _ := http.NewRequest("GET", link, nil)
	request.Header.Set("Host", "i.pximg.net")
	request.Header.Set("Referer", "https://www.pixiv.net/")
	request.Header.Set("Accept", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
	resp, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	// 验证接收到的长度
	length, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	data, _ := ioutil.ReadAll(resp.Body)
	if length != len(data) {
		return nil, "", errors.New("download not complete")
	}
	// 获取文件后缀
	suffix := ".jpg"
	switch resp.Header.Get("Content-Type") {
	case "image/jpeg":
		break
	case "image/png":
		suffix = ".png"
	case "image/gif":
		suffix = ".gif"
	default:
		break
	}
	return &data, suffix, nil
}

// DownloadData 下载 link 到 filedir，返回 filename+suffix, error
func Download(link, filedir, filename string) (string, error) {
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
			return "", err
		}
	}
	data, suffix, err := DownloadData(link)
	if err == nil {
		filepath += suffix
		// 写入文件
		f, _ := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		f.Write(*data)
		f.Close()
		return filepath, nil
	}
	return "", err
}
