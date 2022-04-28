package pixiv

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/process"
	"github.com/FloatTech/zbputils/web"
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

// DownloadToCache 多线程下载第 page 页到 i.Path(page), 返回 error
func (i *Illust) DownloadToCache(page int) error {
	return i.Download(page, i.Path(page))
}

// Download 多线程下载 page 页到 filepath, 返回 error
func (i *Illust) Download(page int, path string) error {
	const slicecap int64 = 65536
	u := i.ImageUrls[page]
	// 获取IP地址
	domain, err := url.Parse(u)
	if err != nil {
		return err
	}

	header := http.Header{
		"Host":          []string{domain.Host},
		"Referer":       []string{"https://www.pixiv.net/"},
		"User-Agent":    []string{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0"},
		"Cache-Control": []string{"no-cache"},
	}

	// 请求 Header
	headreq, err := http.NewRequest("HEAD", u, nil)
	if err != nil {
		return err
	}
	headreq.Header = header.Clone()
	client := web.NewPixivClient()
	headresp, err := client.Do(headreq)
	if err != nil {
		return err
	}

	contentlength, err := strconv.ParseInt(headresp.Header.Get("Content-Length"), 10, 64)
	_ = headresp.Body.Close()
	if err != nil {
		return err
	}

	// 多线程下载
	var wg sync.WaitGroup
	var start int64
	errs := make(chan error, 8)
	buf := make(net.Buffers, 0, contentlength/slicecap+1)
	writers := make([]*binary.Writer, 0, contentlength/slicecap+1)
	index := 0
	for end := math.Min(start+slicecap, contentlength); ; end += slicecap {
		wg.Add(1)
		buf = append(buf, nil)
		writers = append(writers, nil)
		if end > contentlength {
			end = contentlength
		}
		go func(start int64, end int64, index int) {
			// fmt.Println(contentlength, start, end)
			for failedtimes := 0; failedtimes < 3; failedtimes++ {
				req, err := http.NewRequest("GET", u, nil)
				if err != nil {
					errs <- err
					process.SleepAbout1sTo2s()
					continue
				}
				req.Header = header.Clone()
				req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
				resp, err := client.Do(req)
				if err != nil {
					errs <- err
					process.SleepAbout1sTo2s()
					continue
				}
				w := binary.SelectWriter()
				_, err = io.CopyN(w, resp.Body, end-start)
				_ = resp.Body.Close()
				if err != nil {
					errs <- err
					binary.PutWriter(w)
					process.SleepAbout1sTo2s()
					continue
				}
				buf[index] = w.Bytes()
				writers[index] = w
				if err != nil {
					errs <- err
					process.SleepAbout1sTo2s()
					continue
				}
				break
			}
			wg.Done()
		}(start, end, index)
		if end == contentlength {
			break
		}
		start = end
		index++
	}
	msg := ""
	go func() {
		for err := range errs {
			msg += err.Error() + "&"
		}
	}()
	wg.Wait()
	close(errs)
	if msg != "" {
		err = errors.New(msg[:len(msg)-1])
	} else {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, &buf)
		_ = f.Close()
		if err != nil {
			_ = os.Remove(path)
		}
	}
	for _, w := range writers {
		if w != nil {
			binary.PutWriter(w)
		}
	}
	return err
}
