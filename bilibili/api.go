// Package bilibili b站相关API
package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

// ErrAPINeedCookie ...
var ErrAPINeedCookie = errors.New("api need cookie")

// SearchUser 查找b站用户
func SearchUser(cookiecfg *CookieConfig, keyword string) (r []SearchResult, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(SearchUserURL, keyword), nil)
	if err != nil {
		return
	}
	if cookiecfg != nil {
		cookie := ""
		cookie, err = cookiecfg.Load()
		if err != nil {
			return
		}
		req.Header.Add("cookie", cookie)
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("status code: " + strconv.Itoa(res.StatusCode))
		return
	}
	var sd SearchData
	err = json.NewDecoder(res.Body).Decode(&sd)
	if err != nil {
		return
	}
	r = sd.Data.Result
	return
}

// GetVtbDetail 查找vtb信息
func GetVtbDetail(uid string) (result VtbDetail, err error) {
	resp, err := http.Get(fmt.Sprintf(VtbDetailURL, uid))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

// LoadCardDetail 加载卡片
func LoadCardDetail(str string) (card Card, err error) {
	err = json.Unmarshal(binary.StringToBytes(str), &card)
	return
}

// LoadDynamicDetail 加载动态卡片
func LoadDynamicDetail(str string) (card DynamicCard, err error) {
	err = json.Unmarshal(binary.StringToBytes(str), &card)
	return
}

// GetDynamicDetail 用动态id查动态信息
func GetDynamicDetail(cookiecfg *CookieConfig, dynamicIDStr string) (card DynamicCard, err error) {
	var data []byte
	data, err = web.RequestDataWithHeaders(web.NewDefaultClient(), fmt.Sprintf(DynamicDetailURL, dynamicIDStr), "GET", func(req *http.Request) error {
		if cookiecfg != nil {
			cookie := ""
			cookie, err = cookiecfg.Load()
			if err != nil {
				return err
			}
			req.Header.Add("cookie", cookie)
		}
		return nil
	}, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("data.card").Raw), &card)
	return
}

// GetMemberCard 获取b站个人详情
func GetMemberCard(uid any) (result MemberCard, err error) {
	data, err := web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(MemberCardURL, uid), "GET", "", web.RandUA(), nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("data.card").Raw), &result)
	return
}

// GetMedalWall 用b站uid获得牌子
func GetMedalWall(cookiecfg *CookieConfig, uid string) (result []Medal, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(MedalWallURL, uid), nil)
	if err != nil {
		return
	}
	if cookiecfg != nil {
		cookie := ""
		cookie, err = cookiecfg.Load()
		if err != nil {
			return
		}
		req.Header.Add("cookie", cookie)
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	var md MedalData
	err = json.NewDecoder(res.Body).Decode(&md)
	if err != nil {
		return
	}
	if md.Code == -101 {
		err = ErrAPINeedCookie
		return
	}
	if md.Code != 0 {
		err = errors.New(md.Message)
	}
	result = md.Data.List
	return
}

// GetAllGuard 查询mid的上舰信息
func GetAllGuard(mid string) (guardUser GuardUser, err error) {
	var data []byte
	data, err = web.GetData(AllGuardURL)
	if err != nil {
		return
	}
	m := gjson.ParseBytes(data).Get("@this").Map()
	err = json.Unmarshal(binary.StringToBytes(m[mid].String()), &guardUser)
	if err != nil {
		return
	}
	return
}

// GetArticleInfo 用id查专栏信息
func GetArticleInfo(id string) (card Card, err error) {
	var data []byte
	data, err = web.GetData(fmt.Sprintf(ArticleInfoURL, id))
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("data").Raw), &card)
	return
}

// GetLiveRoomInfo 用直播间id查直播间信息
func GetLiveRoomInfo(roomID string) (card RoomCard, err error) {
	var data []byte
	data, err = web.GetData(fmt.Sprintf(ArticleInfoURL, roomID))
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("data").Raw), &card)
	return
}

// GetVideoInfo 用av或bv查视频信息
func GetVideoInfo(id string) (card Card, err error) {
	var data []byte
	_, err = strconv.Atoi(id)
	if err == nil {
		data, err = web.GetData(fmt.Sprintf(VideoInfoURL, id, ""))
	} else {
		data, err = web.GetData(fmt.Sprintf(VideoInfoURL, "", id))
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("data").Raw), &card)
	return
}

// GetVideoSummary 用av或bv查看AI视频总结
func GetVideoSummary(cookiecfg *CookieConfig, id string) (videoSummary VideoSummary, err error) {
	var (
		data []byte
		card Card
	)
	_, err = strconv.Atoi(id)
	if err == nil {
		data, err = web.GetData(fmt.Sprintf(VideoInfoURL, id, ""))
	} else {
		data, err = web.GetData(fmt.Sprintf(VideoInfoURL, "", id))
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("data").Raw), &card)
	if err != nil {
		return
	}
	data, err = web.RequestDataWithHeaders(web.NewDefaultClient(), SignURL(fmt.Sprintf(VideoSummaryURL, card.BvID, card.CID, card.Owner.Mid)), "GET", func(req *http.Request) error {
		if cookiecfg != nil {
			cookie := ""
			cookie, err = cookiecfg.Load()
			if err != nil {
				return err
			}
			req.Header.Add("cookie", cookie)
		}
		req.Header.Set("User-Agent", web.RandUA())
		return nil
	}, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &videoSummary)
	return
}
