package pixiv

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

// P站 无污染 IP 地址
var IPTables = map[string]string{
	"pixiv.net":   "210.140.131.223:443",
	"i.pximg.net": "210.140.92.142:443",
}

//插画结构体
type Illust struct {
	Pid         int64  `db:"pid"`
	Title       string `db:"title"`
	Caption     string `db:"caption"`
	Tags        string `db:"tags"`
	ImageUrls   string `db:"image_urls"`
	AgeLimit    string `db:"age_limit"`
	CreatedTime string `db:"created_time"`
	UserId      int64  `db:"user_id"`
	UserName    string `db:"user_name"`
}

//[]byte转tjson类型
type tjson []byte

//解析json
func (data tjson) Get(path string) gjson.Result {
	return gjson.Get(string(data), path)
}

// Works 获取插画信息
func Works(id int64) (i *Illust, err error) {
	data, err := netPost(fmt.Sprintf("https://pixiv.net/ajax/illust/%d", id))
	if err != nil {
		return nil, err
	}
	json := gjson.ParseBytes(data).Get("body")
	// 如果有"R-18"tag则判断为R-18（暂时）
	var ageLimit = "all-age"
	for _, tag := range json.Get("tags.tags.#.tag").Array() {
		if tag.Str == "R-18" {
			ageLimit = "r18"
			break
		}
	}
	// 解决json返回带html格式
	var caption = strings.ReplaceAll(json.Get("illustComment").Str, "<br />", "\n")
	if index := strings.Index(caption, "<"); index != -1 {
		caption = caption[:index]
	}
	// 解析返回插画信息
	i = &Illust{}
	i.Pid = json.Get("illustId").Int()
	i.Title = json.Get("illustTitle").Str
	i.Caption = caption
	i.Tags = fmt.Sprintln(json.Get("tags.tags.#.tag").Array())
	i.ImageUrls = json.Get("urls.original").Str
	i.AgeLimit = ageLimit
	i.CreatedTime = json.Get("createDate").Str
	i.UserId = json.Get("userId").Int()
	i.UserName = json.Get("userName").Str
	return i, err
}

//搜索元素
type RankValue struct {
	/* required, possible rank modes:
		- daily (default)
	    - weekly
	    - monthly
	    - rookie
	    - original
	    - male
	    - female
	    - daily_r18
	    - weekly_r18
	    - male_r18
	    - female_r18
	    - r18g
	*/
	Mode string
	/* optional, possible rank type:
	    - all (default)
	    - illust
		- ugoira
		- manga
	*/
	Type string
	Page int
	Date string
}

//画作排行榜
func (value RankValue) Rank() (r [18]int, err error) {
	var a []byte
	if value.Mode == "male_r18" || value.Mode == "male" || value.Mode == "female_r18" || value.Mode == "female" {
		value.Type = "all"
		a, err = netPost(fmt.Sprintf("https://pixiv.net/touch/ajax/ranking/illust?mode=%s&type=all&page=%d&date=%s", value.Mode, value.Page, value.Date))
	} else {
		a, err = netPost(fmt.Sprintf("https://pixiv.net/touch/ajax/ranking/illust?mode=%s&type=%s&page=%d&date=%s", value.Mode, value.Type, value.Page, value.Date))
	}
	body := tjson(a)
	for i := 0; i < 18; i++ {
		r[i] = int(body.Get(fmt.Sprintf("body.ranking.%d.illustId", i)).Int())
	}
	return
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
