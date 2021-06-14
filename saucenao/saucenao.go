package saucenao

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type Result struct {
	Similarity float64
	Thumbnail  string
	PixivID    int64
	Title      string
	ImageURL   string
	MemberName string
	MemberID   int64
}

// SauceNaoSearch SauceNao 以图搜图
// 传入图片链接，返回P站结果
func SauceNAO(image string) (*Result, error) {
	var (
		api    = "https://saucenao.com/search.php"
		apiKey = "2cc2772ca550dbacb4c35731a79d341d1a143cb5"

		minSimilarity = 70.0 // 返回图片结果的最小相似度
	)

	// 包装请求参数
	link, _ := url.Parse(api)
	link.RawQuery = url.Values{
		"url":         []string{image},
		"api_key":     []string{apiKey},
		"db":          []string{"5"},
		"numres":      []string{"1"},
		"output_type": []string{"2"},
	}.Encode()

	// 网络请求
	client := &http.Client{}
	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// 如果返回不是200则立刻抛出错误
		return nil, fmt.Errorf("SauceNAO not found, code %d", resp.StatusCode)
	}
	content := gjson.ParseBytes(body)
	if status := content.Get("header.status").Int(); status != 0 {
		// 如果json信息返回status不为0则立刻抛出错误
		return nil, fmt.Errorf("SauceNAO not found, status %d", status)
	}
	if content.Get("results.0.header.similarity").Float() < minSimilarity {
		return nil, fmt.Errorf("SauceNAO not found")
	}
	temp := content.Get("results.0")
	result := &Result{
		Similarity: temp.Get("header.similarity").Float(),
		Thumbnail:  temp.Get("header.thumbnail").Str,
		PixivID:    temp.Get("data.pixiv_id").Int(),
		Title:      temp.Get("data.title").Str,
		ImageURL:   temp.Get("data.ext_urls.0").Str,
		MemberName: temp.Get("data.member_name").Str,
		MemberID:   temp.Get("data.member_id").Int(),
	}
	return result, nil
}
