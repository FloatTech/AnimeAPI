package bilibili

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
	AID             interface{} `json:"aid"`
	BvID            interface{} `json:"bvid"`
	Dynamic         interface{} `json:"dynamic"`
	Pic             string      `json:"pic"`
	Title           string      `json:"title"`
	ID              int         `json:"id"`
	Summary         string      `json:"summary"`
	ImageUrls       []string    `json:"image_urls"`
	OriginImageUrls []string    `json:"origin_image_urls"`
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
	Cover        string      `json:"cover"`
	ShortID      interface{} `json:"short_id"`
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
	Intro      string      `json:"intro"`
	Schema     string      `json:"schema"`
	Author     interface{} `json:"author"`
	AuthorName string      `json:"author_name"`
	PlayCnt    int         `json:"play_cnt"`
	ReplyCnt   int         `json:"reply_cnt"`
	TypeInfo   string      `json:"type_info"`
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
