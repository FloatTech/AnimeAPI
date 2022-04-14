package pixiv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
)

// 插画结构体
type Illust struct {
	Pid         int64    `db:"pid"`
	Title       string   `db:"title"`
	Caption     string   `db:"caption"`
	Tags        string   `db:"tags"`
	ImageUrls   []string `db:"image_urls"`
	AgeLimit    string   `db:"age_limit"`
	CreatedTime string   `db:"created_time"`
	UserId      int64    `db:"user_id"`
	UserName    string   `db:"user_name"`
}

// Works 获取插画信息
func Works(id int64) (i *Illust, err error) {
	data, err := get("https://www.pixiv.net/ajax/illust/" + strconv.FormatInt(id, 10))
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
	u := strings.ReplaceAll(json.Get("urls.original").Str, "_p0.", "_p%d.")
	for j := 0; j < int(json.Get("pageCount").Int()); j++ {
		i.ImageUrls = append(i.ImageUrls, fmt.Sprintf(u, j))
	}
	i.AgeLimit = ageLimit
	i.CreatedTime = json.Get("createDate").Str
	i.UserId = json.Get("userId").Int()
	i.UserName = json.Get("userName").Str
	return i, err
}

// 搜索元素
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

// 画作排行榜
func (value RankValue) Rank() (r [18]int, err error) {
	if value.Mode == "male_r18" || value.Mode == "male" || value.Mode == "female_r18" || value.Mode == "female" {
		value.Type = "all"
	}
	body, err := get(fmt.Sprintf("https://www.pixiv.net/touch/ajax/ranking/illust?mode=%s&type=%s&page=%d&date=%s", value.Mode, value.Type, value.Page, value.Date))
	i := 0
	gjson.Get(binary.BytesToString(body), "body.ranking").ForEach(func(key, value gjson.Result) bool {
		r[i] = int(value.Get("illustId").Int())
		i++
		if i == 18 {
			return false
		}
		return true
	})
	return
}

// get 返回请求数据
func get(link string) ([]byte, error) {
	return web.RequestDataWith(
		web.NewPixivClient(), link, "GET",
		"https://www.pixiv.net/",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0",
	)
}
