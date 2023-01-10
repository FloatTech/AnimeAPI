package qzone

// EmotionPublishRequest 发说说请求体
type EmotionPublishRequest struct {
	CodeVersion    string `json:"code_version"`
	Con            string `json:"con"`
	Feedversion    string `json:"feedversion"`
	Format         string `json:"format"`
	Hostuin        string `json:"hostuin"`
	Paramstr       string `json:"paramstr"`
	PicBo          string `json:"pic_bo"`
	PicTemplate    string `json:"pic_template"`
	Qzreferrer     string `json:"qzreferrer"`
	Richtype       string `json:"richtype"`
	Richval        string `json:"richval"`
	SpecialURL     string `json:"special_url"`
	Subrichtype    string `json:"subrichtype"`
	SynTweetVerson string `json:"syn_tweet_verson"`
	ToSign         string `json:"to_sign"`
	UgcRight       string `json:"ugc_right"`
	Ver            string `json:"ver"`
	Who            string `json:"who"`
}

// EmotionPublishVo 发说说响应体
type EmotionPublishVo struct {
	Activity     []interface{} `json:"activity"`
	Attach       interface{}   `json:"attach"`
	AuthFlag     int           `json:"auth_flag"`
	Code         int           `json:"code"`
	Conlist      []Conlist     `json:"conlist"`
	Content      string        `json:"content"`
	Message      string        `json:"message"`
	OurlInfo     interface{}   `json:"ourl_info"`
	PicTemplate  string        `json:"pic_template"`
	Right        int           `json:"right"`
	Secret       int           `json:"secret"`
	Signin       int           `json:"signin"`
	Smoothpolicy Smoothpolicy1 `json:"smoothpolicy"`
	Subcode      int           `json:"subcode"`
	T1Icon       int           `json:"t1_icon"`
	T1Name       string        `json:"t1_name"`
	T1Ntime      int           `json:"t1_ntime"`
	T1Source     int           `json:"t1_source"`
	T1SourceName string        `json:"t1_source_name"`
	T1SourceURL  string        `json:"t1_source_url"`
	T1Tid        string        `json:"t1_tid"`
	T1Time       string        `json:"t1_time"`
	T1Uin        int           `json:"t1_uin"`
	ToTweet      int           `json:"to_tweet"`
	UgcRight     int           `json:"ugc_right"`
}

// Conlist 说说文字消息
type Conlist struct {
	Con  string `json:"con"`
	Type int    `json:"type"`
}

// Smoothpolicy1 暂定
type Smoothpolicy1 struct {
	Smoothpolicy Smoothpolicy
}

// Smoothpolicy 暂定
type Smoothpolicy struct {
	ComswDisableSosoSearch  int `json:"comsw.disable_soso_search"`
	L1SwReadFirstCacheOnly  int `json:"l1sw.read_first_cache_only"`
	L2SwDontGetReplyCmt     int `json:"l2sw.dont_get_reply_cmt"`
	L2SwMixsvrFrdnumPerTime int `json:"l2sw.mixsvr_frdnum_per_time"`
	L3SwHideReplyCmt        int `json:"l3sw.hide_reply_cmt"`
	L4SwReadTdbOnly         int `json:"l4sw.read_tdb_only"`
	L5SwReadCacheOnly       int `json:"l5sw.read_cache_only"`
}

// UploadImageRequest 上传图片请求体
type UploadImageRequest struct {
	Albumtype        string `json:"albumtype"`
	BackUrls         string `json:"backUrls"`
	Base64           string `json:"base64"`
	Charset          string `json:"charset"`
	Exttype          string `json:"exttype"`
	Filename         string `json:"filename"`
	HdHeight         string `json:"hd_height"`
	HdQuality        string `json:"hd_quality"`
	HdWidth          string `json:"hd_width"`
	JsonhtmlCallback string `json:"jsonhtml_callback"`
	OutputCharset    string `json:"output_charset"`
	OutputType       string `json:"output_type"`
	PSkey            string `json:"p_skey"`
	PUin             string `json:"p_uin"`
	Picfile          string `json:"picfile"`
	Qzonetoken       string `json:"qzonetoken"`
	Qzreferrer       string `json:"qzreferrer"`
	Refer            string `json:"refer"`
	Skey             string `json:"skey"`
	Uin              string `json:"uin"`
	UploadHd         string `json:"upload_hd"`
	Uploadtype       string `json:"uploadtype"`
	URL              string `json:"url"`
	Zzpanelkey       string `json:"zzpanelkey"`
	Zzpaneluin       string `json:"zzpaneluin"`
}

// UploadImageVo 上传图片响应体
type UploadImageVo struct {
	Data struct {
		Pre          string `json:"pre"`
		URL          string `json:"url"`
		Lloc         string `json:"lloc"`
		Sloc         string `json:"sloc"`
		Type         int    `json:"type"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		Albumid      string `json:"albumid"`
		Totalpic     int    `json:"totalpic"`
		Limitpic     int    `json:"limitpic"`
		OriginURL    string `json:"origin_url"`
		OriginUUID   string `json:"origin_uuid"`
		OriginWidth  int    `json:"origin_width"`
		OriginHeight int    `json:"origin_height"`
		Contentlen   int    `json:"contentlen"`
	} `json:"data"`
	Ret int `json:"ret"`
}

// MsgListRequest 说说列表请求体
type MsgListRequest struct {
	Callback           string `json:"callback"`
	CodeVersion        string `json:"code_version"`
	Format             string `json:"format"`
	Ftype              string `json:"ftype"`
	GTk                string `json:"g_tk"`
	NeedPrivateComment string `json:"need_private_comment"`
	Num                string `json:"num"`
	Pos                string `json:"pos"`
	Replynum           string `json:"replynum"`
	Sort               string `json:"sort"`
	Uin                string `json:"uin"`
}

// MsgListVo 说说列表响应体
type MsgListVo struct {
	AuthFlag     int          `json:"auth_flag"`
	CensorCount  int          `json:"censor_count"`
	CensorFlag   int          `json:"censor_flag"`
	CensorTotal  int          `json:"censor_total"`
	Cginame      int          `json:"cginame"`
	Code         int          `json:"code"`
	Logininfo    Logininfo    `json:"logininfo"`
	Mentioncount int          `json:"mentioncount"`
	Message      string       `json:"message"`
	Msglist      []Msglist    `json:"msglist"`
	Name         string       `json:"name"`
	Num          int          `json:"num"`
	Sign         int          `json:"sign"`
	Smoothpolicy Smoothpolicy `json:"smoothpolicy"`
	Subcode      int          `json:"subcode"`
	Timertotal   int          `json:"timertotal"`
	Total        int          `json:"total"`
	Usrinfo      Usrinfo      `json:"usrinfo"`
}

// Logininfo 登录信息
type Logininfo struct {
	Name string `json:"name"`
	Uin  int    `json:"uin"`
}

// Lbs 位置信息
type Lbs struct {
	ID     string `json:"id"`
	Idname string `json:"idname"`
	Name   string `json:"name"`
	PosX   string `json:"pos_x"`
	PosY   string `json:"pos_y"`
}

// Pic 图片信息
type Pic struct {
	AbsolutePosition int    `json:"absolute_position"`
	BHeight          int    `json:"b_height"`
	BWidth           int    `json:"b_width"`
	Curlikekey       string `json:"curlikekey"`
	Height           int    `json:"height"`
	PicID            string `json:"pic_id"`
	Pictype          int    `json:"pictype"`
	Richsubtype      int    `json:"richsubtype"`
	Rtype            int    `json:"rtype"`
	Smallurl         string `json:"smallurl"`
	Unilikekey       string `json:"unilikekey"`
	URL1             string `json:"url1"`
	URL2             string `json:"url2"`
	URL3             string `json:"url3"`
	Who              int    `json:"who"`
	Width            int    `json:"width"`
}

// Msglist 单个说说的详细信息
type Msglist struct {
	Certified   int       `json:"certified"`
	Cmtnum      int       `json:"cmtnum"`
	Conlist     []Conlist `json:"conlist"`
	Content     string    `json:"content"`
	CreateTime  string    `json:"createTime"`
	CreatedTime int       `json:"created_time"`
	EditMask    int64     `json:"editMask"`
	Fwdnum      int       `json:"fwdnum"`
	HasMoreCon  int       `json:"has_more_con"`
	IsEditable  int       `json:"isEditable"`
	Issigin     int       `json:"issigin"`
	Lastmodify  int       `json:"lastmodify"`
	Lbs         Lbs       `json:"lbs"`
	Name        string    `json:"name"`
	PicTemplate string    `json:"pic_template"`
	Right       int       `json:"right"`
	RtSum       int       `json:"rt_sum"`
	Secret      int       `json:"secret"`
	SourceAppid string    `json:"source_appid"`
	SourceName  string    `json:"source_name"`
	SourceURL   string    `json:"source_url"`
	T1Source    int       `json:"t1_source"`
	T1Subtype   int       `json:"t1_subtype"`
	T1Termtype  int       `json:"t1_termtype"`
	Tid         string    `json:"tid"`
	UgcRight    int       `json:"ugc_right"`
	Uin         int       `json:"uin"`
	Wbid        int       `json:"wbid"`
	Pic         []Pic     `json:"pic,omitempty"`
	Pictotal    int       `json:"pictotal,omitempty"`
}

// Usrinfo 个人信息
type Usrinfo struct {
	Concern    int    `json:"concern"`
	CreateTime string `json:"createTime"`
	Fans       int    `json:"fans"`
	Followed   int    `json:"followed"`
	Msg        string `json:"msg"`
	Msgnum     int    `json:"msgnum"`
	Name       string `json:"name"`
	Uin        int    `json:"uin"`
}

// LikeRequest 空间点赞请求体
type LikeRequest struct {
	Curkey     string `json:"curkey"`
	Face       string `json:"face"`
	From       string `json:"from"`
	Fupdate    string `json:"fupdate"`
	Opuin      string `json:"opuin"`
	Qzreferrer string `json:"qzreferrer"`
	Unikey     string `json:"unikey"`
	Format     string `json:"format"`
}
