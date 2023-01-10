// Package yandex yandex搜图
package yandex

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/FloatTech/AnimeAPI/pixiv"
	xpath "github.com/antchfx/htmlquery"
)

// Yandex yandex搜图
func Yandex(image string) (*pixiv.Illust, error) {
	search, _ := url.Parse("https://yandex.com/images/search")
	search.RawQuery = url.Values{
		"rpt":  []string{"imageview"},
		"url":  []string{image},
		"site": []string{"pixiv.net"},
	}.Encode()
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
	}

	// 网络请求
	req, _ := http.NewRequest("GET", search.String(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := xpath.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	// 取出每个返回的结果
	list := xpath.Find(doc, `/html/body/div[3]/div[2]/div[1]/div/div[1]/div[1]/div[2]/div[2]/div/section/div[2]/div[1]/div[1]/a`)
	if len(list) != 1 {
		return nil, errors.New("Yandex not found")
	}
	link := list[0].Attr[1].Val
	dest, _ := url.Parse(link)
	rawid := dest.Query().Get("illust_id")
	if rawid == "" {
		return nil, errors.New("Yandex not found")
	}
	// 链接取出PIXIV id
	id, _ := strconv.ParseInt(rawid, 10, 64)
	if id == 0 {
		return nil, errors.New("convert to pid error")
	}

	illust, err := pixiv.Works(id)
	if err != nil {
		return nil, err
	}
	// 待完善
	return illust, nil
}
