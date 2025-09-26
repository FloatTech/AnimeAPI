package bilibili

import (
	"net/http"
	"strconv"
	"time"

	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	msgType = map[int]string{
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

// HumanNum 格式化人数
func HumanNum(res int) string {
	if res/10000 != 0 {
		return strconv.FormatFloat(float64(res)/10000, 'f', 2, 64) + "万"
	}
	return strconv.Itoa(res)
}

// GetRealURL 获取跳转后的链接
func GetRealURL(url string) (realurl string, err error) {
	data, err := http.Head(url)
	if err != nil {
		return
	}
	_ = data.Body.Close()
	realurl = data.Request.URL.String()
	return
}

// card2msg cType=1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308时,处理Card字符串,cType为card类型
func card2msg(dynamicCard *DynamicCard, card *Card, cType int) (msg []message.Segment, err error) {
	var (
		vote Vote
	)
	msg = make([]message.Segment, 0, 16)
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
