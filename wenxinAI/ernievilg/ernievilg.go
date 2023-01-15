// Package ernievilg 百度文心AI画图
package ernievilg

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/web"
)

type tokendata struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// GetToken 获取当天的token
//
// 申请账号链接:https://wenxin.baidu.com/moduleApi/key
//
// clientID为API key,clientSecret为Secret key
//
// token有效时间为24小时
func GetToken(clientID, clientSecret string) (token string, err error) {
	requestURL := "https://wenxin.baidu.com/moduleApi/portal/api/oauth/token?grant_type=client_credentials&client_id=" + url.QueryEscape(clientID) + "&client_secret=" + url.QueryEscape(clientSecret)
	postData := url.Values{}
	postData.Add("name", "ATRI")
	postData.Add("language", "golang")
	data, err := web.PostData(requestURL, "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}
	var parsed tokendata
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return
	}
	if parsed.Msg != "success" {
		err = errors.New(parsed.Msg + ",code:" + strconv.Itoa(parsed.Code))
		return
	}
	return parsed.Data, nil
}

type workstate struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		RequestID string `json:"requestId"`
		TaskID    int    `json:"taskId"`
	} `json:"data"`
}

// BuildWork 创建画图任务
//
// token:GetToken函数获取,
//
// keyword:图片描述,长度不超过64个字,prompt指南:https://wenxin.baidu.com/wenxin/docs#Ol7ece95m
//
// picType:图片风格，目前支持风格有：油画、水彩画、卡通、粉笔画、儿童画、蜡笔画
//
// picSize:图片尺寸，目前支持的有：1024*1024 方图、1024*1536 长图、1536*1024 横图。
// 传入的是尺寸数值，非文字。
//
// taskID:任务ID，用于查询结果
func BuildWork(token, keyword, picType, picSize string) (taskID int, err error) {
	requestURL := "https://wenxin.baidu.com/moduleApi/portal/api/rest/1.0/ernievilg/v1/txt2img?access_token=" + url.QueryEscape(token) +
		"&text=" + url.QueryEscape(keyword) + "&style=" + url.QueryEscape(picType) + "&resolution=" + url.QueryEscape(picSize)
	postData := url.Values{}
	postData.Add("name", "ATRI")
	postData.Add("language", "golang")
	data, err := web.PostData(requestURL, "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}
	var parsed workstate
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return
	}
	if parsed.Msg != "success" {
		err = errors.New(parsed.Msg + ",code:" + strconv.Itoa(parsed.Code))
		return
	}
	return parsed.Data.TaskID, nil
}

type picdata struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Img        string   `json:"img"`
		Waiting    string   `json:"waiting"`
		ImgUrls    []PicURL `json:"imgUrls"`
		CreateTime string   `json:"createTime"`
		RequestID  string   `json:"requestId"`
		Style      string   `json:"style"`
		Text       string   `json:"text"`
		Resolution string   `json:"resolution"`
		TaskID     int      `json:"taskId"`
		Status     int      `json:"status"`
	} `json:"data"`
}

// PicURL ...
type PicURL struct {
	Image string      `json:"image"`
	Score interface{} `json:"score"`
}

// GetPic 获取图片内容
//
// token由GetToken函数获取,taskID由BuildWork函数获取
//
// PicURL:[x]struct{Image:图片链接,Score:评分}
//
// API会返回x张图片,数量不确定的,随机的。
//
// 评分目前都是null,我不知道有什么用，既然API预留了，我也预留吧
//
// stauts:结果状态,"30s"代表还在排队生成,"0"表示结果OK
func GetPic(token string, taskID int) (picurls []PicURL, status string, err error) {
	requestURL := "https://wenxin.baidu.com/moduleApi/portal/api/rest/1.0/ernievilg/v1/getImg?access_token=" + url.QueryEscape(token) +
		"&taskId=" + strconv.Itoa(taskID)
	postData := url.Values{}
	postData.Add("name", "ATRI")
	postData.Add("language", "golang")
	data, err := web.PostData(requestURL, "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}
	var parsed picdata
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return
	}
	if parsed.Msg != "success" {
		err = errors.New(parsed.Msg + ",code:" + strconv.Itoa(parsed.Code))
		return
	}
	status = parsed.Data.Waiting
	picurls = parsed.Data.ImgUrls
	return
}
