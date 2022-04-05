package nsfw

import (
	"bytes"
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

const apiurl = "https://nsfwtag.azurewebsites.net/api/nsfw?url="

func Classify(u string) (*Picture, error) {
	u = apiurl + url.QueryEscape(u)
	var data []byte
	data, err := web.GetData(u)
	if err != nil {
		return nil, err
	}
	ps := make([]Picture, 1)
	err = json.Unmarshal(bytes.ReplaceAll(data, []byte("'"), []byte("\"")), &ps)
	if err != nil {
		return nil, err
	}
	return &ps[0], nil
}
