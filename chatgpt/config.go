package chatgpt

import "time"

const API = "https://chat.openai.com/"

type Config struct {
	UA              string
	SessionToken    string
	CFClearance     string
	RefreshInterval time.Duration
	Timeout         time.Duration
}
