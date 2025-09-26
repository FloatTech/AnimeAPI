package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// TURL bilibili动态前缀
	TURL = "https://t.bilibili.com/"
	// LiveURL bilibili直播前缀
	LiveURL = "https://live.bilibili.com/"
	// DynamicDetailURL 当前动态信息,一个card
	DynamicDetailURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=%v"
	// MemberCardURL 个人信息
	MemberCardURL = "https://api.bilibili.com/x/web-interface/card?mid=%v"
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
	VideoSummaryURL = "https://api.bilibili.com/x/web-interface/view/conclusion/get?bvid=%v&cid=%v&up_mid=%v"
	// VideoDownloadURL 视频下载
	VideoDownloadURL = "https://api.bilibili.com/x/player/playurl?bvid=%v&cid=%v&qn=80&fnval=1&fnver=0&fourk=1"
	// OnlineTotalURL 在线人数
	OnlineTotalURL = "https://api.bilibili.com/x/player/online/total?bvid=%v&cid=%v"
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

// ToMessage 动态转消息
func (dynamicCard *DynamicCard) ToMessage() (msg []message.Segment, err error) {
	var (
		card  Card
		vote  Vote
		cType int
	)
	msg = make([]message.Segment, 0, 16)
	// 初始化结构体
	err = json.Unmarshal(binary.StringToBytes(dynamicCard.Card), &card)
	if err != nil {
		return
	}
	if dynamicCard.Extension.Vote != "" {
		err = json.Unmarshal(binary.StringToBytes(dynamicCard.Extension.Vote), &vote)
		if err != nil {
			return
		}
	}
	cType = dynamicCard.Desc.Type
	// 生成消息
	switch cType {
	case 1:
		msg = append(msg, message.Text(card.User.Uname, msgType[cType], "\n",
			card.Item.Content, "\n",
			"转发的内容: \n"))
		var originMsg []message.Segment
		var co Card
		co, err = LoadCardDetail(card.Origin)
		if err != nil {
			return
		}
		originMsg, err = card2msg(dynamicCard, &co, card.Item.OrigType)
		if err != nil {
			return
		}
		msg = append(msg, originMsg...)
	case 2:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.UploadTime), 0).Format("2006-01-02 15:04:05"), msgType[cType], "\n",
			card.Item.Description))
		for i := 0; i < len(card.Item.Pictures); i++ {
			msg = append(msg, message.Image(card.Item.Pictures[i].ImgSrc))
		}
	case 4:
		msg = append(msg, message.Text(card.User.Uname, "在", time.Unix(int64(card.Item.Timestamp), 0).Format("2006-01-02 15:04:05"), msgType[cType], "\n",
			card.Item.Content, "\n"))
		if dynamicCard.Extension.Vote != "" {
			msg = append(msg, message.Text("【投票】", vote.Desc, "\n",
				"截止日期: ", time.Unix(int64(vote.Endtime), 0).Format("2006-01-02 15:04:05"), "\n",
				"参与人数: ", HumanNum(vote.JoinNum), "\n",
				"投票选项( 最多选择", vote.ChoiceCnt, "项 )\n"))
			for i := 0; i < len(vote.Options); i++ {
				msg = append(msg, message.Text("- ", vote.Options[i].Idx, ". ", vote.Options[i].Desc, "\n"))
				if vote.Options[i].ImgURL != "" {
					msg = append(msg, message.Image(vote.Options[i].ImgURL))
				}
			}
		}
	case 8:
		msg = append(msg, message.Text(card.Owner.Name, "在", time.Unix(int64(card.Pubdate), 0).Format("2006-01-02 15:04:05"), msgType[cType], "\n",
			card.Title))
		msg = append(msg, message.Image(card.Pic))
		msg = append(msg, message.Text(card.Desc, "\n",
			card.ShareSubtitle, "\n",
			"视频链接: ", card.ShortLink, "\n"))
	case 16:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.UploadTime), 0).Format("2006-01-02 15:04:05"), msgType[cType], "\n",
			card.Item.Description))
		msg = append(msg, message.Image(card.Item.Cover.Default))
	case 64:
		msg = append(msg, message.Text(card.Author.(map[string]any)["name"], "在", time.Unix(int64(card.PublishTime), 0).Format("2006-01-02 15:04:05"), msgType[cType], "\n",
			card.Title, "\n",
			card.Summary))
		for i := 0; i < len(card.ImageUrls); i++ {
			msg = append(msg, message.Image(card.ImageUrls[i]))
		}
		if card.ID != 0 {
			msg = append(msg, message.Text("文章链接: https://www.bilibili.com/read/cv", card.ID, "\n"))
		}
	case 256:
		msg = append(msg, message.Text(card.Upper, "在", time.Unix(int64(card.Ctime), 0).Format("2006-01-02 15:04:05"), msgType[cType], "\n",
			card.Title))
		msg = append(msg, message.Image(card.Cover))
		msg = append(msg, message.Text(card.Intro, "\n"))
		if card.ID != 0 {
			msg = append(msg, message.Text("音频链接: https://www.bilibili.com/audio/au", card.ID, "\n"))
		}

	case 2048:
		msg = append(msg, message.Text(card.User.Uname, msgType[cType], "\n",
			card.Vest.Content, "\n",
			card.Sketch.Title, "\n",
			card.Sketch.DescText, "\n"))
		msg = append(msg, message.Image(card.Sketch.CoverURL))
		msg = append(msg, message.Text("分享链接: ", card.Sketch.TargetURL, "\n"))
	case 4308:
		if dynamicCard.Desc.UserProfile.Info.Uname != "" {
			msg = append(msg, message.Text(dynamicCard.Desc.UserProfile.Info.Uname, msgType[cType], "\n"))
		}
		msg = append(msg, message.Image(card.LivePlayInfo.Cover))
		msg = append(msg, message.Text("\n", card.LivePlayInfo.Title, "\n",
			"房间号: ", card.LivePlayInfo.RoomID, "\n",
			"分区: ", card.LivePlayInfo.ParentAreaName))
		if card.LivePlayInfo.ParentAreaName != card.LivePlayInfo.AreaName {
			msg = append(msg, message.Text("-", card.LivePlayInfo.AreaName))
		}
		if card.LivePlayInfo.LiveStatus == 0 {
			msg = append(msg, message.Text("未开播 \n"))
		} else {
			msg = append(msg, message.Text("直播中 ", card.LivePlayInfo.WatchedShow, "\n"))
		}
		msg = append(msg, message.Text("直播链接: ", card.LivePlayInfo.Link))
	default:
		msg = append(msg, message.Text("动态id: ", dynamicCard.Desc.DynamicIDStr, "未知动态类型: ", cType, "\n"))
	}
	if dynamicCard.Desc.DynamicIDStr != "" {
		msg = append(msg, message.Text("动态链接: ", TURL, dynamicCard.Desc.DynamicIDStr))
	}
	return
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

// ToArticleMessage 专栏转消息
func (card *Card) ToArticleMessage(defaultID string) (msg []message.Segment) {
	msg = make([]message.Segment, 0, len(card.OriginImageUrls)+1)
	for i := 0; i < len(card.OriginImageUrls); i++ {
		msg = append(msg, message.Image(card.OriginImageUrls[i]))
	}
	msg = append(msg, message.Text("\n", card.Title, "\n", "UP主: ", card.AuthorName, "\n",
		"阅读: ", HumanNum(card.Stats.View), " 评论: ", HumanNum(card.Stats.Reply), "\n",
		CVURL, defaultID))
	return
}

// ToVideoMessage 视频卡片转消息
func (card *Card) ToVideoMessage() (msg []message.Segment, err error) {
	var (
		mCard       MemberCard
		onlineTotal OnlineTotal
	)
	msg = make([]message.Segment, 0, 16)
	mCard, err = GetMemberCard(card.Owner.Mid)
	msg = append(msg, message.Text("标题: ", card.Title, "\n"))
	if card.Rights.IsCooperation == 1 {
		for i := 0; i < len(card.Staff); i++ {
			msg = append(msg, message.Text(card.Staff[i].Title, ": ", card.Staff[i].Name, " 粉丝: ", HumanNum(card.Staff[i].Follower), "\n"))
		}
	} else {
		if err != nil {
			msg = append(msg, message.Text("UP主: ", card.Owner.Name, "\n"))
		} else {
			msg = append(msg, message.Text("UP主: ", card.Owner.Name, " 粉丝: ", HumanNum(mCard.Fans), "\n"))
		}
	}
	msg = append(msg, message.Image(card.Pic))
	data, err := web.GetData(fmt.Sprintf(OnlineTotalURL, card.BvID, card.CID))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &onlineTotal)
	if err != nil {
		return
	}
	msg = append(msg, message.Text("👀播放: ", HumanNum(card.Stat.View), " 💬弹幕: ", HumanNum(card.Stat.Danmaku),
		"\n👍点赞: ", HumanNum(card.Stat.Like), " 💰投币: ", HumanNum(card.Stat.Coin),
		"\n📁收藏: ", HumanNum(card.Stat.Favorite), " 🔗分享: ", HumanNum(card.Stat.Share),
		"\n📝简介: ", card.Desc,
		"\n🏄‍♂️ 总共 ", onlineTotal.Data.Total, " 人在观看，", onlineTotal.Data.Count, " 人在网页端观看\n",
		VURL, card.BvID, "\n\n"))
	return
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

// ToMessage 直播卡片转消息
func (card *RoomCard) ToMessage() (msg []message.Segment) {
	msg = make([]message.Segment, 0, 10)
	msg = append(msg, message.Image(card.RoomInfo.Keyframe))
	msg = append(msg, message.Text("\n", card.RoomInfo.Title, "\n",
		"主播: ", card.AnchorInfo.BaseInfo.Uname, "\n",
		"房间号: ", card.RoomInfo.RoomID, "\n"))
	if card.RoomInfo.ShortID != 0 {
		msg = append(msg, message.Text("短号: ", card.RoomInfo.ShortID, "\n"))
	}
	msg = append(msg, message.Text("分区: ", card.RoomInfo.ParentAreaName))
	if card.RoomInfo.ParentAreaName != card.RoomInfo.AreaName {
		msg = append(msg, message.Text("-", card.RoomInfo.AreaName))
	}
	if card.RoomInfo.LiveStatus == 0 {
		msg = append(msg, message.Text("未开播 \n"))
	} else {
		msg = append(msg, message.Text("直播中 ", HumanNum(card.RoomInfo.Online), "人气\n"))
	}
	if card.RoomInfo.ShortID != 0 {
		msg = append(msg, message.Text("直播间链接: ", LURL, card.RoomInfo.ShortID))
	} else {
		msg = append(msg, message.Text("直播间链接: ", LURL, card.RoomInfo.RoomID))
	}

	return
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
					LastLiveIncome       float64       `json:"lastLiveIncome"`
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

// VideoDownload 视频下载结构体(mp4格式)
type VideoDownload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		From              string   `json:"from"`
		Result            string   `json:"result"`
		Message           string   `json:"message"`
		Quality           int      `json:"quality"`
		Format            string   `json:"format"`
		Timelength        int      `json:"timelength"`
		AcceptFormat      string   `json:"accept_format"`
		AcceptDescription []string `json:"accept_description"`
		AcceptQuality     []int    `json:"accept_quality"`
		VideoCodecid      int      `json:"video_codecid"`
		SeekParam         string   `json:"seek_param"`
		SeekType          string   `json:"seek_type"`
		Durl              []struct {
			Order     int      `json:"order"`
			Length    int      `json:"length"`
			Size      int      `json:"size"`
			Ahead     string   `json:"ahead"`
			Vhead     string   `json:"vhead"`
			URL       string   `json:"url"`
			BackupURL []string `json:"backup_url"`
		} `json:"durl"`
		SupportFormats []struct {
			Quality        int         `json:"quality"`
			Format         string      `json:"format"`
			NewDescription string      `json:"new_description"`
			DisplayDesc    string      `json:"display_desc"`
			Superscript    string      `json:"superscript"`
			Codecs         interface{} `json:"codecs"`
		} `json:"support_formats"`
		HighFormat   interface{} `json:"high_format"`
		LastPlayTime int         `json:"last_play_time"`
		LastPlayCid  int         `json:"last_play_cid"`
	} `json:"data"`
}

// OnlineTotal 在线人数结构体
type OnlineTotal struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Total      string `json:"total"`
		Count      string `json:"count"`
		ShowSwitch struct {
			Total bool `json:"total"`
			Count bool `json:"count"`
		} `json:"show_switch"`
		Abtest struct {
			Group string `json:"group"`
		} `json:"abtest"`
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
