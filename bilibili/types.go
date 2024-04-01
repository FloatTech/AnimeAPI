package bilibili

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/FloatTech/floatbox/file"
)

const (
	// TURL bilibili动态前缀
	TURL = "https://t.bilibili.com/"
	// LiveURL bilibili直播前缀
	LiveURL = "https://live.bilibili.com/"
	// DynamicDetailURL 当前动态信息,一个card
	DynamicDetailURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=%v"
	// MemberCardURL 个人信息
	MemberCardURL = "https://account.bilibili.com/api/member/getCardByMid?mid=%v"
	// ArticleInfoURL 查看专栏信息
	ArticleInfoURL = "https://api.bilibili.com/x/article/viewinfo?id=%v"
	// CVURL b站专栏前缀
	CVURL = "https://www.bilibili.com/read/cv"
	// LiveRoomInfoURL 查看直播间信息
	LiveRoomInfoURL = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom?room_id=%v"
	// LURL b站直播间前缀
	LURL = "https://live.bilibili.com/"
	// VideoInfoURL 查看视频信息
	VideoInfoURL = "https://api.bilibili.com/x/web-interface/view?aid=%v&bvid=%v"
	// VURL 视频网址前缀
	VURL = "https://www.bilibili.com/video/"
	// SearchUserURL 查找b站用户
	SearchUserURL = "http://api.bilibili.com/x/web-interface/search/type?search_type=bili_user&keyword=%v"
	// VtbDetailURL 查找vtb信息
	VtbDetailURL = "https://api.vtbs.moe/v1/detail/%v"
	// MedalWallURL 查找牌子
	MedalWallURL = "https://api.live.bilibili.com/xlive/web-ucenter/user/MedalWall?target_id=%v"
	// SpaceHistoryURL 历史动态信息,一共12个card
	SpaceHistoryURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?host_uid=%v&offset_dynamic_id=%v&need_top=0"
	// LiveListURL 获得直播状态
	LiveListURL = "https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids"
	// DanmakuAPI 弹幕网获得用户弹幕api
	DanmakuAPI = "https://ukamnads.icu/api/v2/user?uId=%v&pageNum=%v&pageSize=5&target=-1&useEmoji=true"
	// DanmakuURL 弹幕网链接
	DanmakuURL = "https://danmakus.com/user/%v"
	// AllGuardURL 查询所有舰长,提督,总督
	AllGuardURL = "https://api.vtbs.moe/v1/guard/all"
	// VideoSummaryURL AI视频总结
	VideoSummaryURL = "https://api.bilibili.com/x/web-interface/view/conclusion/get?bvid=%v&cid=%v"
	// NavURL 导航URL
	NavURL = "https://api.bilibili.com/x/web-interface/nav"
)

// DynamicCard 总动态结构体,包括desc,card
type DynamicCard struct {
	Desc      Desc   `json:"desc"`
	Card      string `json:"card"`
	Extension struct {
		VoteCfg struct {
			VoteID  int    `json:"vote_id"`
			Desc    string `json:"desc"`
			JoinNum int    `json:"join_num"`
		} `json:"vote_cfg"`
		Vote string `json:"vote"`
	} `json:"extension"`
}

// Card 卡片结构体
type Card struct {
	Item struct {
		Content     string `json:"content"`
		UploadTime  int    `json:"upload_time"`
		Description string `json:"description"`
		Pictures    []struct {
			ImgSrc string `json:"img_src"`
		} `json:"pictures"`
		Timestamp int `json:"timestamp"`
		Cover     struct {
			Default string `json:"default"`
		} `json:"cover"`
		OrigType int `json:"orig_type"`
	} `json:"item"`
	AID             any      `json:"aid"`
	BvID            any      `json:"bvid"`
	Dynamic         any      `json:"dynamic"`
	CID             int      `json:"cid"`
	Pic             string   `json:"pic"`
	Title           string   `json:"title"`
	ID              int      `json:"id"`
	Summary         string   `json:"summary"`
	ImageUrls       []string `json:"image_urls"`
	OriginImageUrls []string `json:"origin_image_urls"`
	Sketch          struct {
		Title     string `json:"title"`
		DescText  string `json:"desc_text"`
		CoverURL  string `json:"cover_url"`
		TargetURL string `json:"target_url"`
	} `json:"sketch"`
	Stat struct {
		Aid      int `json:"aid"`
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		Share    int `json:"share"`
		Like     int `json:"like"`
	} `json:"stat"`
	Stats struct {
		Aid      int `json:"aid"`
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		Share    int `json:"share"`
		Like     int `json:"like"`
	} `json:"stats"`
	Owner struct {
		Name    string `json:"name"`
		Pubdate int    `json:"pubdate"`
		Mid     int    `json:"mid"`
	} `json:"owner"`
	Cover        string `json:"cover"`
	ShortID      any    `json:"short_id"`
	LivePlayInfo struct {
		ParentAreaName string `json:"parent_area_name"`
		AreaName       string `json:"area_name"`
		Cover          string `json:"cover"`
		Link           string `json:"link"`
		Online         int    `json:"online"`
		RoomID         int    `json:"room_id"`
		LiveStatus     int    `json:"live_status"`
		WatchedShow    string `json:"watched_show"`
		Title          string `json:"title"`
	} `json:"live_play_info"`
	Intro      string `json:"intro"`
	Schema     string `json:"schema"`
	Author     any    `json:"author"`
	AuthorName string `json:"author_name"`
	PlayCnt    int    `json:"play_cnt"`
	ReplyCnt   int    `json:"reply_cnt"`
	TypeInfo   string `json:"type_info"`
	User       struct {
		Name  string `json:"name"`
		Uname string `json:"uname"`
	} `json:"user"`
	Desc          string `json:"desc"`
	ShareSubtitle string `json:"share_subtitle"`
	ShortLink     string `json:"short_link"`
	PublishTime   int    `json:"publish_time"`
	BannerURL     string `json:"banner_url"`
	Ctime         int    `json:"ctime"`
	Vest          struct {
		Content string `json:"content"`
	} `json:"vest"`
	Upper   string `json:"upper"`
	Origin  string `json:"origin"`
	Pubdate int    `json:"pubdate"`
	Rights  struct {
		IsCooperation int `json:"is_cooperation"`
	} `json:"rights"`
	Staff []struct {
		Title    string `json:"title"`
		Name     string `json:"name"`
		Follower int    `json:"follower"`
	} `json:"staff"`
}

// Desc 描述结构体
type Desc struct {
	Type         int    `json:"type"`
	DynamicIDStr string `json:"dynamic_id_str"`
	OrigType     int    `json:"orig_type"`
	Timestamp    int    `json:"timestamp"`
	Origin       struct {
		DynamicIDStr string `json:"dynamic_id_str"`
	} `json:"origin"`
	UserProfile struct {
		Info struct {
			Uname string `json:"uname"`
		} `json:"info"`
	} `json:"user_profile"`
}

// Vote 投票结构体
type Vote struct {
	ChoiceCnt int    `json:"choice_cnt"`
	Desc      string `json:"desc"`
	Endtime   int    `json:"endtime"`
	JoinNum   int    `json:"join_num"`
	Options   []struct {
		Idx    int    `json:"idx"`
		Desc   string `json:"desc"`
		ImgURL string `json:"img_url"`
	} `json:"options"`
}

// MemberCard 个人信息卡片
type MemberCard struct {
	Mid        string  `json:"mid"`
	Name       string  `json:"name"`
	Sex        string  `json:"sex"`
	Face       string  `json:"face"`
	Coins      float64 `json:"coins"`
	Regtime    int64   `json:"regtime"`
	Birthday   string  `json:"birthday"`
	Sign       string  `json:"sign"`
	Attentions []int64 `json:"attentions"`
	Fans       int     `json:"fans"`
	Friend     int     `json:"friend"`
	Attention  int     `json:"attention"`
	LevelInfo  struct {
		CurrentLevel int `json:"current_level"`
	} `json:"level_info"`
}

// RoomCard 直播间卡片
type RoomCard struct {
	RoomInfo struct {
		RoomID         int    `json:"room_id"`
		ShortID        int    `json:"short_id"`
		Title          string `json:"title"`
		LiveStatus     int    `json:"live_status"`
		AreaName       string `json:"area_name"`
		ParentAreaName string `json:"parent_area_name"`
		Keyframe       string `json:"keyframe"`
		Online         int    `json:"online"`
	} `json:"room_info"`
	AnchorInfo struct {
		BaseInfo struct {
			Uname string `json:"uname"`
		} `json:"base_info"`
	} `json:"anchor_info"`
}

// SearchData 查找b站用户总结构体
type SearchData struct {
	Data struct {
		NumResults int            `json:"numResults"`
		Result     []SearchResult `json:"result"`
	} `json:"data"`
}

// SearchResult 查找b站用户结果
type SearchResult struct {
	Mid    int64  `json:"mid"`
	Uname  string `json:"uname"`
	Gender int64  `json:"gender"`
	Usign  string `json:"usign"`
	Level  int64  `json:"level"`
}

// MedalData 牌子接口返回结构体
type MedalData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []Medal `json:"list"`
	} `json:"data"`
}

// MedalInfo b站牌子信息
type MedalInfo struct {
	Mid              int64  `json:"target_id"`
	MedalName        string `json:"medal_name"`
	Level            int64  `json:"level"`
	MedalColorStart  int64  `json:"medal_color_start"`
	MedalColorEnd    int64  `json:"medal_color_end"`
	MedalColorBorder int64  `json:"medal_color_border"`
}

// Medal ...
type Medal struct {
	Uname     string `json:"target_name"`
	MedalInfo `json:"medal_info"`
}

// MedalSorter ...
type MedalSorter []Medal

// Len ...
func (m MedalSorter) Len() int {
	return len(m)
}

// Swap ...
func (m MedalSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less ...
func (m MedalSorter) Less(i, j int) bool {
	return m[i].Level > m[j].Level
}

// VtbDetail vtb信息
type VtbDetail struct {
	Mid      int    `json:"mid"`
	Uname    string `json:"uname"`
	Video    int    `json:"video"`
	Roomid   int    `json:"roomid"`
	Rise     int    `json:"rise"`
	Follower int    `json:"follower"`
	GuardNum int    `json:"guardNum"`
	AreaRank int    `json:"areaRank"`
}

// GuardUser dd用户
type GuardUser struct {
	Uname string    `json:"uname"`
	Face  string    `json:"face"`
	Mid   int64     `json:"mid"`
	Dd    [][]int64 `json:"dd"`
}

// Danmakusuki 弹幕网结构体
type Danmakusuki struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total    int  `json:"total"`
		PageNum  int  `json:"pageNum"`
		PageSize int  `json:"pageSize"`
		HasMore  bool `json:"hasMore"`
		Data     struct {
			Records []struct {
				Channel struct {
					UID                  int           `json:"uId"`
					UName                string        `json:"uName"`
					RoomID               int           `json:"roomId"`
					FaceURL              string        `json:"faceUrl"`
					FrameURL             string        `json:"frameUrl"`
					IsLiving             bool          `json:"isLiving"`
					Title                string        `json:"title"`
					Tags                 []interface{} `json:"tags"`
					LastLiveDate         int64         `json:"lastLiveDate"`
					LastLiveDanmakuCount int           `json:"lastLiveDanmakuCount"`
					TotalDanmakuCount    int           `json:"totalDanmakuCount"`
					TotalIncome          float64       `json:"totalIncome"`
					TotalLiveCount       int           `json:"totalLiveCount"`
					TotalLiveSecond      int           `json:"totalLiveSecond"`
					AddDate              string        `json:"addDate"`
					CommentCount         int           `json:"commentCount"`
					LastLiveIncome       int           `json:"lastLiveIncome"`
				} `json:"channel"`
				Live struct {
					LiveID           string  `json:"liveId"`
					IsFinish         bool    `json:"isFinish"`
					IsFull           bool    `json:"isFull"`
					ParentArea       string  `json:"parentArea"`
					Area             string  `json:"area"`
					CoverURL         string  `json:"coverUrl"`
					DanmakusCount    int     `json:"danmakusCount"`
					StartDate        int64   `json:"startDate"`
					StopDate         int64   `json:"stopDate"`
					Title            string  `json:"title"`
					TotalIncome      float64 `json:"totalIncome"`
					WatchCount       int     `json:"watchCount"`
					LikeCount        int     `json:"likeCount"`
					PayCount         int     `json:"payCount"`
					InteractionCount int     `json:"interactionCount"`
					MaxOnlineCount   int     `json:"maxOnlineCount"`
				} `json:"live"`
				Danmakus []struct {
					UID      int     `json:"uId"`
					UName    string  `json:"uName"`
					Type     int64   `json:"type"`
					SendDate int64   `json:"sendDate"`
					Message  string  `json:"message"`
					Price    float64 `json:"price"`
				} `json:"danmakus"`
			} `json:"records"`
		} `json:"data"`
	} `json:"data"`
}

// VideoSummary AI视频总结结构体
type VideoSummary struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Code        int `json:"code"`
		ModelResult struct {
			ResultType int    `json:"result_type"`
			Summary    string `json:"summary"`
			Outline    []struct {
				Title       string `json:"title"`
				PartOutline []struct {
					Timestamp int    `json:"timestamp"`
					Content   string `json:"content"`
				} `json:"part_outline"`
				Timestamp int `json:"timestamp"`
			} `json:"outline"`
		} `json:"model_result"`
		Stid       string `json:"stid"`
		Status     int    `json:"status"`
		LikeNum    int    `json:"like_num"`
		DislikeNum int    `json:"dislike_num"`
	} `json:"data"`
}

// CookieConfig 配置结构体
type CookieConfig struct {
	BilibiliCookie string `json:"bilibili_cookie"`
	file           string
}

// NewCookieConfig ...
func NewCookieConfig(file string) *CookieConfig {
	return &CookieConfig{
		file: file,
	}
}

// Set ...
func (cfg *CookieConfig) Set(cookie string) (err error) {
	cfg.BilibiliCookie = cookie
	return cfg.Save()
}

// Load ...
func (cfg *CookieConfig) Load() (cookie string, err error) {
	if cfg.BilibiliCookie != "" {
		cookie = cfg.BilibiliCookie
		return
	}
	if file.IsNotExist(cfg.file) {
		err = errors.New("no cookie config")
		return
	}
	reader, err := os.Open(cfg.file)
	if err != nil {
		return
	}
	defer reader.Close()
	err = json.NewDecoder(reader).Decode(cfg)
	cookie = cfg.BilibiliCookie
	return
}

// Save ...
func (cfg *CookieConfig) Save() (err error) {
	reader, err := os.Create(cfg.file)
	if err != nil {
		return err
	}
	defer reader.Close()
	return json.NewEncoder(reader).Encode(cfg)
}
