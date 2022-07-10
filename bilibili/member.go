package bilibili

const (
	// MemberCardURL 个人信息
	MemberCardURL = "https://account.bilibili.com/api/member/getCardByMid?mid=%v"
)

// MemberCard 个人信息卡片
type MemberCard struct {
	Mid        string  `json:"mid"`
	Name       string  `json:"name"`
	Sex        string  `json:"sex"`
	Face       string  `json:"face"`
	Coins      float64 `json:"coins"`
	Regtime    int     `json:"regtime"`
	Birthday   string  `json:"birthday"`
	Sign       string  `json:"sign"`
	Attentions []int   `json:"attentions"`
	Fans       int     `json:"fans"`
	Friend     int     `json:"friend"`
	Attention  int     `json:"attention"`
	LevelInfo  struct {
		CurrentLevel int `json:"current_level"`
	} `json:"level_info"`
}
