package nsfw

import (
	"encoding/json"
	"net/url"

	"github.com/FloatTech/zbputils/web"
)

type Picture struct {
	Sexy     float64 `json:"sexy"`
	Neutral  float64 `json:"neutral"`
	Porn     float64 `json:"porn"`
	Hentai   float64 `json:"hentai"`
	Drawings float64 `json:"drawings"`
}

const apiurl = "https://sayuri.fumiama.top/nsfw?"

func Classify(urls ...string) (p []Picture, err error) {
	u := apiurl
	for _, s := range urls {
		u += "urls=" + url.QueryEscape(s) + "&"
	}
	u = u[:len(u)-1]
	var data []byte
	data, err = web.GetData(u)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &p)
	return
}
