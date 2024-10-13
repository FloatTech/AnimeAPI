// Package emozi 颜文字抽象转写
package emozi

const api = "https://emozi.seku.su/api/"

// User 用户
type User struct {
	name string
	pswd string
	auth string
}

type encodebody struct {
	Random bool   `json:"random"`
	Text   string `json:"text"`
	Choice []int  `json:"choice"`
}

type encoderesult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Text   string `json:"text"`
		Choice []int  `json:"choice,omitempty"`
	} `json:"result"`
}

type decodebody struct {
	Force bool   `json:"force"`
	Text  string `json:"text"`
}

type decoderesult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type loginbody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
