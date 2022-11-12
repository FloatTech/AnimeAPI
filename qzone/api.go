// Package qzone QQ空间API
package qzone

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/FloatTech/floatbox/binary"
)

var (
	cRe = regexp.MustCompile(`_Callback\((.*)\)`)
)

const (
	userQzoneURL      = "https://user.qzone.qq.com"
	ua                = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	contentType       = "application/x-www-form-urlencoded"
	params            = "g_tk=%v"
	inpcqqURL         = "https://h5.qzone.qq.com/feeds/inpcqq?uin=%v&qqver=5749&timestamp=%v"
	emotionPublishURL = userQzoneURL + "/proxy/domain/taotao.qzone.qq.com/cgi-bin/emotion_cgi_publish_v6?" + params
	uploadImageURL    = "https://up.qzone.qq.com/cgi-bin/upload/cgi_upload_image?" + params
	msglistURL        = userQzoneURL + "/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msglist_v6"
	likeURL           = userQzoneURL + "/proxy/domain/w.qzone.qq.com/cgi-bin/likes/internal_dolike_app?" + params
	ptqrshowURL       = "https://ssl.ptlogin2.qq.com/ptqrshow?appid=549000912&e=2&l=M&s=3&d=72&v=4&t=0.31232733520361844&daid=5&pt_3rd_aid=0"
	ptqrloginURL      = "https://xui.ptlogin2.qq.com/ssl/ptqrlogin?u1=https://qzs.qq.com/qzone/v5/loginsucc.html?para=izone&ptqrtoken=%v&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-1656992258324&js_ver=22070111&js_type=1&login_sig=&pt_uistyle=40&aid=549000912&daid=5&has_onekey=1&&o1vId=1e61428d61cb5015701ad73d5fb59f73"
	checkSigURL       = "https://ptlogin2.qzone.qq.com/check_sig?pttype=1&uin=%v&service=ptqrlogin&nodirect=1&ptsigx=%v&s_url=https://qzs.qq.com/qzone/v5/loginsucc.html?para=izone&f_url=&ptlang=2052&ptredirect=100&aid=549000912&daid=5&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=3&pt_aid=0&pt_aaid=16&pt_light=0&pt_3rd_aid=0"
)

// Ptqrshow 获得登录二维码
func Ptqrshow() (data []byte, qrsig string, ptqrtoken string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	ptqrshowReq, err := http.NewRequest("GET", ptqrshowURL, nil)
	if err != nil {
		return
	}
	ptqrshowResp, err := client.Do(ptqrshowReq)
	if err != nil {
		return
	}
	defer ptqrshowResp.Body.Close()
	for _, v := range ptqrshowResp.Cookies() {
		if v.Name == "qrsig" {
			qrsig = v.Value
			break
		}
	}
	if qrsig == "" {
		return
	}
	ptqrtoken = genderGTK(qrsig, 0)
	data, err = io.ReadAll(ptqrshowResp.Body)
	return
}

// Ptqrlogin 登录回调
func Ptqrlogin(qrsig string, qrtoken string) (data []byte, cookie string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	ptqrloginReq, err := http.NewRequest("GET", fmt.Sprintf(ptqrloginURL, qrtoken), nil)
	if err != nil {
		return
	}
	ptqrloginReq.Header.Add("cookie", "qrsig="+qrsig)
	ptqrloginResp, err := client.Do(ptqrloginReq)
	if err != nil {
		return
	}
	defer ptqrloginResp.Body.Close()
	for _, v := range ptqrloginReq.Cookies() {
		if v.Value != "" {
			cookie += v.Name + "=" + v.Value + ";"
		}
	}
	data, err = io.ReadAll(ptqrloginResp.Body)
	return
}

// LoginRedirect 登录成功回调
func LoginRedirect(redirectURL string) (cookie string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	u, err := url.Parse(redirectURL)
	if err != nil {
		return
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}
	redirectReq, err := http.NewRequest("GET", fmt.Sprintf(checkSigURL, values["uin"][0], values["ptsigx"][0]), nil)
	if err != nil {
		return
	}
	redirectResp, err := client.Do(redirectReq)
	if err != nil {
		return
	}
	defer redirectResp.Body.Close()
	for _, v := range redirectResp.Cookies() {
		if v.Value != "" {
			cookie += v.Name + "=" + v.Value + ";"
		}
	}
	return
}

// Manager qq空间信息管理
type Manager struct {
	Cookie string
	QQ     string
	Gtk    string
	Gtk2   string
	PSkey  string
	Skey   string
	Uin    string
}

// NewManager 初始化信息
func NewManager(cookie string) (m Manager) {
	cookie = strings.ReplaceAll(cookie, " ", "")
	for _, v := range strings.Split(cookie, ";") {
		name, val, f := strings.Cut(v, "=")
		if f {
			switch name {
			case "uin":
				m.Uin = val
			case "skey":
				m.Skey = val
			case "p_skey":
				m.PSkey = val
			}
		}
	}
	m.Gtk = genderGTK(m.Skey, 5381)
	m.Gtk2 = genderGTK(m.PSkey, 5381)
	m.QQ = strings.TrimPrefix(m.Uin, "o")
	m.Cookie = cookie
	return
}

// EmotionPublishRaw 发送说说
func (m *Manager) EmotionPublishRaw(epr EmotionPublishRequest) (result EmotionPublishVo, err error) {
	client := &http.Client{}
	payload := strings.NewReader(structToStr(epr))
	request, err := http.NewRequest("POST", fmt.Sprintf(emotionPublishURL, m.Gtk2), payload)
	if err != nil {
		return
	}
	request.Header.Add("referer", userQzoneURL)
	request.Header.Add("origin", userQzoneURL)
	request.Header.Add("cookie", m.Cookie)
	request.Header.Add("user-agent", ua)
	request.Header.Add("content-type", contentType)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&result)
	return
}

// UploadImage 上传图片
func (m *Manager) UploadImage(base64img string) (result UploadImageVo, err error) {
	uir := UploadImageRequest{
		Filename:      "filename",
		Uin:           m.QQ,
		Skey:          m.Skey,
		Zzpaneluin:    m.QQ,
		PUin:          m.QQ,
		PSkey:         m.PSkey,
		Uploadtype:    "1",
		Albumtype:     "7",
		Exttype:       "0",
		Refer:         "shuoshuo",
		OutputType:    "json",
		Charset:       "utf-8",
		OutputCharset: "utf-8",
		UploadHd:      "1",
		HdWidth:       "2048",
		HdHeight:      "10000",
		HdQuality:     "96",
		BackUrls:      "http://upbak.photo.qzone.qq.com/cgi-bin/upload/cgi_upload_image,http://119.147.64.75/cgi-bin/upload/cgi_upload_image",
		URL:           fmt.Sprintf(uploadImageURL, m.Gtk2),
		Base64:        "1",
		Picfile:       base64img,
		Qzreferrer:    userQzoneURL + "/" + m.QQ,
	}

	payload := strings.NewReader(structToStr(uir))
	client := &http.Client{}
	request, err := http.NewRequest("POST", fmt.Sprintf(uploadImageURL, m.Gtk2), payload)
	if err != nil {
		return
	}
	request.Header.Add("referer", userQzoneURL)
	request.Header.Add("origin", userQzoneURL)
	request.Header.Add("cookie", m.Cookie)
	request.Header.Add("user-agent", ua)
	request.Header.Add("content-type", contentType)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	r := cRe.FindStringSubmatch(binary.BytesToString(data))
	if len(r) < 2 {
		err = errors.New("上传失败")
		return
	}
	err = json.Unmarshal(binary.StringToBytes(r[1]), &result)
	return
}

// EmotionPublish 发送说说,content是文字,base64imgList是base64图片
func (m *Manager) EmotionPublish(content string, base64imgList []string) (result EmotionPublishVo, err error) {

	var (
		uir         UploadImageVo
		picBo       string
		richval     string
		richvalList = make([]string, 0, 9)
		picBoList   = make([]string, 0, 9)
	)

	for _, base64img := range base64imgList {
		uir, err = m.UploadImage(base64img)
		if err != nil {
			return
		}
		picBo, richval, err = getPicBoAndRichval(uir)
		if err != nil {
			return
		}
		richvalList = append(richvalList, richval)
		picBoList = append(picBoList, picBo)
	}

	epr := EmotionPublishRequest{
		SynTweetVerson: "1",
		Paramstr:       "1",
		Who:            "1",
		Con:            content,
		Feedversion:    "1",
		Ver:            "1",
		UgcRight:       "1",
		ToSign:         "0",
		Hostuin:        m.QQ,
		CodeVersion:    "1",
		Format:         "json",
		Qzreferrer:     userQzoneURL + "/" + m.QQ,
	}
	if len(base64imgList) > 0 {
		epr.Richtype = "1"
		epr.Richval = strings.Join(richvalList, "\t")
		epr.Subrichtype = "1"
		epr.PicBo = strings.Join(picBoList, ",")
	}

	result, err = m.EmotionPublishRaw(epr)
	return
}

// EmotionMsglist 获取说说列表
func (m *Manager) EmotionMsglist(num string, replynum string) (mlv MsgListVo, err error) {
	mlr := MsgListRequest{
		Uin:                m.QQ,
		Ftype:              "0",
		Sort:               "0",
		Num:                num,
		Replynum:           replynum,
		GTk:                m.Gtk2,
		Callback:           "_preloadCallback",
		CodeVersion:        "1",
		Format:             "json",
		NeedPrivateComment: "1",
	}
	mlv, err = m.EmotionMsglistRaw(mlr)
	return
}

// EmotionMsglistRaw 获取说说列表
func (m *Manager) EmotionMsglistRaw(mlr MsgListRequest) (mlv MsgListVo, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", msglistURL+"?"+structToStr(mlr), nil)
	if err != nil {
		return
	}
	request.Header.Add("referer", userQzoneURL)
	request.Header.Add("origin", userQzoneURL)
	request.Header.Add("cookie", m.Cookie)
	request.Header.Add("user-agent", ua)
	request.Header.Add("content-type", contentType)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&mlv)
	return
}

// LikeRaw 空间点赞(貌似只能给自己点赞,预留)
func (m *Manager) LikeRaw(lr LikeRequest) (err error) {
	client := &http.Client{}
	payload := strings.NewReader(structToStr(lr))
	request, err := http.NewRequest("POST", likeURL, payload)
	if err != nil {
		return
	}
	request.Header.Add("referer", userQzoneURL)
	request.Header.Add("origin", userQzoneURL)
	request.Header.Add("cookie", m.Cookie)
	request.Header.Add("user-agent", ua)
	request.Header.Add("content-type", contentType)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}
	fmt.Printf("data:%s\n", data)
	return
}

func getPicBoAndRichval(data UploadImageVo) (picBo, richval string, err error) {
	var flag bool
	if data.Ret != 0 {
		err = errors.New("上传失败")
		return
	}
	_, picBo, flag = strings.Cut(data.Data.URL, "&bo=")
	if !flag {
		err = errors.New("上传图片返回的地址错误")
		return
	}
	richval = fmt.Sprintf(",%s,%s,%s,%d,%d,%d,,%d,%d", data.Data.Albumid, data.Data.Lloc, data.Data.Sloc, data.Data.Type, data.Data.Height, data.Data.Width, data.Data.Height, data.Data.Width)
	return
}
