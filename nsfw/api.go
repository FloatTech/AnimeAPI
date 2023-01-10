// Package nsfw 图片鉴赏
package nsfw

import (
	"encoding/json"
	"net/url"

	"github.com/FloatTech/floatbox/web"
)

// Picture ...
type Picture struct {
	Sexy     float64 `json:"sexy"`
	Neutral  float64 `json:"neutral"`
	Porn     float64 `json:"porn"`
	Hentai   float64 `json:"hentai"`
	Drawings float64 `json:"drawings"`
}

const apiurl = "https://nsfwtag.azurewebsites.net/api/nsfw?url="

// Classify ...
func Classify(u string) (*Picture, error) {
	u = apiurl + url.QueryEscape(u)
	var data []byte
	data, err := web.GetData(u)
	if err != nil {
		return nil, err
	}
	ps := make([]Picture, 1)
	err = json.Unmarshal(data, &ps)
	if err != nil {
		return nil, err
	}
	return &ps[0], nil
}
