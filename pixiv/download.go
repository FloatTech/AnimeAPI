package pixiv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
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

// DownAll 单线程下载 link 到 filedir，返回 filedir+filename+suffix, error
func DownAll(link, filedir, filename string) (string, error) {
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

// Download 多线程下载 link 到 filedir，返回 filedir+filename+suffix, error
func Download(link, filedir, filename string) (string, error) {
	var slicecap int64 = 65536
	// 获取IP地址
	domain, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	// P站特殊客户端
	client := &http.Client{
		// 解决中国大陆无法访问的问题
		Transport: &http.Transport{
			// 更改 dns
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp", IPTables[domain.Host])
			},
			// 隐藏 sni 标志
			TLSClientConfig: &tls.Config{
				ServerName:         "-",
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
		},
	}
	header := http.Header{
		"Host":       []string{domain.Host},
		"Referer":    []string{"https://www.pixiv.net/"},
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0"},
	}

	// 请求 Header
	headreq, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
		return "", err
	}
	headreq.Header = header.Clone()
	headresp, err := client.Do(headreq)
	if err != nil {
		return "", err
	}
	defer headresp.Body.Close()

	contentlength, _ := strconv.ParseInt(headresp.Header.Get("Content-Length"), 10, 64)
	var suffix string
	switch headresp.Header.Get("Content-Type") {
	case "image/jpeg":
		suffix = ".jpg"
	case "image/png":
		suffix = ".png"
	case "image/gif":
		suffix = ".gif"
	default:
		suffix = ".jpg"
	}
	var filepath = filedir + filename + suffix
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	// 多线程下载
	wg := sync.WaitGroup{}
	var start int64
	for {
		end := func(a int64, b int64) int64 {
			if a > b {
				return b
			}
			return a
		}(start+slicecap, contentlength)
		wg.Add(1)
		go func(f *os.File, start int64, end int64) {
			var failedtimes int
			fmt.Println(contentlength, start, end)
			for {
				if failedtimes >= 3 {
					break
				}
				req, err := http.NewRequest("GET", link, nil)
				if err != nil {
					failedtimes++
					continue
				}
				req.Header = header.Clone()
				req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
				resp, err := client.Do(req)
				if err != nil {
					failedtimes++
					continue
				}
				defer resp.Body.Close()
				b, _ := ioutil.ReadAll(resp.Body)
				if len(b) != int(end-start) {
					failedtimes++
					continue
				}
				_, err = f.WriteAt(b, int64(start))
				if err != nil {
					failedtimes++
					continue
				}
				break
			}
			wg.Done()
		}(f, start, end)
		start = end
		if start >= contentlength {
			break
		}
	}
	wg.Wait()
	return filepath, nil
}
