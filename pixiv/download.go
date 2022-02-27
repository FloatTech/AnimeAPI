package pixiv

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/FloatTech/zbputils/math"
)

const CacheDir = "data/pixiv/"

func init() {
	err := os.MkdirAll(CacheDir, 0755)
	if err != nil {
		panic(err)
	}
}

// Path 图片本地缓存路径
func (i *Illust) Path(page int) string {
	u := i.ImageUrls[page]
	f := CacheDir + u[strings.LastIndex(u, "/")+1:]
	return f
}

// DownloadToCache 多线程下载第 page 页到 i.Path(page)，返回 error
func (i *Illust) DownloadToCache(page int) error {
	f := i.Path(page)
	file, err := os.Create(f)
	if err != nil {
		return err
	}
	err = i.Download(page, file)
	_ = file.Sync()
	stat, err1 := file.Stat()
	var size int64
	if err1 != nil {
		size = stat.Size()
	}
	_ = file.Close()
	if err != nil || size <= 0 {
		_ = os.Remove(f)
	}
	return err
}

// Download 多线程下载 link 到 filepath，返回 error
func (i *Illust) Download(page int, f *os.File) error {
	const slicecap int64 = 65536
	u := i.ImageUrls[page]
	// 获取IP地址
	domain, err := url.Parse(u)
	if err != nil {
		return err
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
	headreq, err := http.NewRequest("HEAD", u, nil)
	if err != nil {
		return err
	}
	headreq.Header = header.Clone()
	headresp, err := client.Do(headreq)
	if err != nil {
		return err
	}
	defer headresp.Body.Close()

	contentlength, _ := strconv.ParseInt(headresp.Header.Get("Content-Length"), 10, 64)
	// 多线程下载
	wg := sync.WaitGroup{}
	var start int64
	var mu sync.Mutex
	for {
		end := math.Min64(start+slicecap, contentlength)
		wg.Add(1)
		go func(start int64, end int64) {
			var failedtimes int
			// fmt.Println(contentlength, start, end)
			for {
				if failedtimes >= 3 {
					break
				}
				req, err := http.NewRequest("GET", u, nil)
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
				mu.Lock()
				_, err = f.WriteAt(b, int64(start))
				mu.Unlock()
				if err != nil {
					failedtimes++
					continue
				}
				break
			}
			wg.Done()
		}(start, end)
		if end >= contentlength {
			break
		}
		start = end
	}
	wg.Wait()
	return nil
}
