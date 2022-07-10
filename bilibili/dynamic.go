// Package bilibili b站动态解析库
package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// TURL bilibili动态前缀
	TURL = "https://t.bilibili.com/"
	// SpaceHistoryURL 历史动态信息,一共12个card
	SpaceHistoryURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?host_uid=%v&offset_dynamic_id=%v&need_top=0"
	// DynamicDetailURL 当前动态信息,一个card
	DynamicDetailURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=%v"
)

var (
	typeMsg = map[int]string{
		1:    "转发了动态",
		2:    "有图营业",
		4:    "无图营业",
		8:    "投稿了视频",
		16:   "投稿了短视频",
		64:   "投稿了文章",
		256:  "投稿了音频",
		2048: "发布了简报",
		4200: "发布了直播",
		4308: "发布了直播",
	}
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

// DynamicCard2msg cType=0时,处理DynCard字符串,cType=1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308时,处理Card字符串,cType为card类型
func DynamicCard2msg(str string, cType int) (msg []message.MessageSegment, err error) {
	var (
		dynamicCard DynamicCard
		card        Card
		vote        Vote
	)
	msg = make([]message.MessageSegment, 0, 16)
	// 初始化结构体
	switch cType {
	case 0:
		err = json.Unmarshal(binary.StringToBytes(str), &dynamicCard)
		if err != nil {
			return
		}
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
	case 1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308:
		err = json.Unmarshal(binary.StringToBytes(str), &card)
		if err != nil {
			return
		}
	default:
		err = errors.New("只有0, 1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308模式")
		return
	}
	// 生成消息
	switch cType {
	case 1:
		msg = append(msg, message.Text(card.User.Uname, typeMsg[cType], "\n",
			card.Item.Content, "\n",
			"转发的内容: \n"))
		var originMsg []message.MessageSegment
		originMsg, err = DynamicCard2msg(card.Origin, card.Item.OrigType)
		if err != nil {
			return
		}
		msg = append(msg, originMsg...)
	case 2:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.UploadTime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n",
			card.Item.Description))
		for i := 0; i < len(card.Item.Pictures); i++ {
			msg = append(msg, message.Image(card.Item.Pictures[i].ImgSrc))
		}
	case 4:
		msg = append(msg, message.Text(card.User.Uname, "在", time.Unix(int64(card.Item.Timestamp), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n",
			card.Item.Content, "\n"))
		if dynamicCard.Extension.Vote != "" {
			msg = append(msg, message.Text("【投票】", vote.Desc, "\n",
				"截止日期: ", time.Unix(int64(vote.Endtime), 0).Format("2006-01-02 15:04:05"), "\n",
				"参与人数: ", humanNum(vote.JoinNum), "\n",
				"投票选项( 最多选择", vote.ChoiceCnt, "项 )\n"))
			for i := 0; i < len(vote.Options); i++ {
				msg = append(msg, message.Text("- ", vote.Options[i].Idx, ". ", vote.Options[i].Desc, "\n"))
				if vote.Options[i].ImgURL != "" {
					msg = append(msg, message.Image(vote.Options[i].ImgURL))
				}
			}
		}
	case 8:
		msg = append(msg, message.Text(card.Owner.Name, "在", time.Unix(int64(card.Pubdate), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n",
			card.Title))
		msg = append(msg, message.Image(card.Pic))
		msg = append(msg, message.Text(card.Desc, "\n",
			card.ShareSubtitle, "\n",
			"视频链接: ", card.ShortLink, "\n"))
	case 16:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.UploadTime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n",
			card.Item.Description))
		msg = append(msg, message.Image(card.Item.Cover.Default))
	case 64:
		msg = append(msg, message.Text(card.Author.(map[string]interface{})["name"], "在", time.Unix(int64(card.PublishTime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n",
			card.Title, "\n",
			card.Summary))
		for i := 0; i < len(card.ImageUrls); i++ {
			msg = append(msg, message.Image(card.ImageUrls[i]))
		}
		if card.ID != 0 {
			msg = append(msg, message.Text("文章链接: https://www.bilibili.com/read/cv", card.ID, "\n"))
		}
	case 256:
		msg = append(msg, message.Text(card.Upper, "在", time.Unix(int64(card.Ctime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n",
			card.Title))
		msg = append(msg, message.Image(card.Cover))
		msg = append(msg, message.Text(card.Intro, "\n"))
		if card.ID != 0 {
			msg = append(msg, message.Text("音频链接: https://www.bilibili.com/audio/au", card.ID, "\n"))
		}

	case 2048:
		msg = append(msg, message.Text(card.User.Uname, typeMsg[cType], "\n",
			card.Vest.Content, "\n",
			card.Sketch.Title, "\n",
			card.Sketch.DescText, "\n"))
		msg = append(msg, message.Image(card.Sketch.CoverURL))
		msg = append(msg, message.Text("分享链接: ", card.Sketch.TargetURL, "\n"))
	case 4308:
		if dynamicCard.Desc.UserProfile.Info.Uname != "" {
			msg = append(msg, message.Text(dynamicCard.Desc.UserProfile.Info.Uname, typeMsg[cType], "\n"))
		}
		msg = append(msg, message.Image(card.LivePlayInfo.Cover))
		msg = append(msg, message.Text(card.LivePlayInfo.Title, "\n",
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

// DynamicDetail 用动态id查动态信息
func DynamicDetail(dynamicIDStr string) (msg []message.MessageSegment, err error) {
	var data []byte
	data, err = web.GetData(fmt.Sprintf(DynamicDetailURL, dynamicIDStr))
	if err != nil {
		return
	}
	return DynamicCard2msg(gjson.ParseBytes(data).Get("data.card").Raw, 0)
}
