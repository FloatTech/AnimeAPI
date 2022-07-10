package bilibili

import (
	"encoding/json"
	"fmt"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	// LiveRoomInfoURL 查看直播间信息
	LiveRoomInfoURL = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom?room_id=%v"
	// UserLiveURL 查看直播用户信息
	UserLiveURL = "https://api.bilibili.com/x/space/acc/info?mid=%v"
	// LURL b站直播间前缀
	LURL = "https://live.bilibili.com/"
)

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

// LiveCard2msg 直播卡片转消息
func LiveCard2msg(str string) (msg []message.MessageSegment, err error) {
	var (
		card RoomCard
	)
	msg = make([]message.MessageSegment, 0, 16)
	err = json.Unmarshal(binary.StringToBytes(str), &card)
	if err != nil {
		return
	}
	msg = append(msg, message.Image(card.RoomInfo.Keyframe))
	msg = append(msg, message.Text(card.RoomInfo.Title, "\n",
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
		msg = append(msg, message.Text("直播中 ", humanNum(card.RoomInfo.Online), "人气\n"))
	}
	if card.RoomInfo.ShortID != 0 {
		msg = append(msg, message.Text("直播间链接: ", LURL, card.RoomInfo.ShortID))
	} else {
		msg = append(msg, message.Text("直播间链接: ", LURL, card.RoomInfo.RoomID))
	}

	return
}

// LiveRoomInfo 用直播间id查直播间信息
func LiveRoomInfo(roomID string) (msg []message.MessageSegment, err error) {
	var data []byte
	data, err = web.GetData(fmt.Sprintf(LiveRoomInfoURL, roomID))
	if err != nil {
		return
	}
	return LiveCard2msg(gjson.ParseBytes(data).Get("data").Raw)
}
