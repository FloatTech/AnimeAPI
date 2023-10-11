package bilibili

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/RomiChan/syncx"
	"github.com/tidwall/gjson"
)

var (
	mixinKeyEncTab = []int{
		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}
	cache          syncx.Map[string, string]
	lastUpdateTime time.Time
	replacements   = [...]string{"!", "'", "(", ")", "*"}
)

// SignURL wbi签名包装 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/misc/sign/wbi.md
func SignURL(urlStr string) string {
	urlObj, _ := url.Parse(urlStr)
	imgKey, subKey := getWbiKeysCached()
	query := urlObj.Query()
	params := map[string]string{}
	for k, v := range query {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	newParams := wbiSign(params, imgKey, subKey)
	for k, v := range newParams {
		query.Set(k, v)
	}
	urlObj.RawQuery = query.Encode()
	newURL := urlObj.String()
	return newURL
}

func getMixinKey(orig string) string {
	var str strings.Builder
	t := 0
	for _, v := range mixinKeyEncTab {
		if v < len(orig) {
			str.WriteByte(orig[v])
			t++
		}
		if t > 31 {
			break
		}
	}
	return str.String()
}

func wbiSign(params map[string]string, imgKey string, subKey string) map[string]string {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	params["wts"] = currTime
	// Sort keys
	keys := make([]string, 0, len(params))
	for k, v := range params {
		keys = append(keys, k)
		for _, old := range replacements {
			v = strings.ReplaceAll(v, old, "")
		}
		params[k] = v
	}
	sort.Strings(keys)
	h := md5.New()
	for k, v := range keys {
		h.Write([]byte(v))
		h.Write([]byte{'='})
		h.Write([]byte(params[v]))
		if k < len(keys)-1 {
			h.Write([]byte{'&'})
		}
	}
	h.Write([]byte(mixinKey))
	params["w_rid"] = hex.EncodeToString(h.Sum(make([]byte, 0, md5.Size)))
	return params
}

func getWbiKeysCached() (string, string) {
	if time.Since(lastUpdateTime).Minutes() > 10 {
		imgKey, subKey := getWbiKeys()
		cache.Store("imgKey", imgKey)
		cache.Store("subKey", subKey)
		lastUpdateTime = time.Now()
		return imgKey, subKey
	}
	imgKeyI, _ := cache.Load("imgKey")
	subKeyI, _ := cache.Load("subKey")
	return imgKeyI, subKeyI
}

func getWbiKeys() (string, string) {
	data, _ := web.GetData(NavURL)
	json := binary.BytesToString(data)
	imgURL := gjson.Get(json, "data.wbi_img.img_url").String()
	subURL := gjson.Get(json, "data.wbi_img.sub_url").String()
	imgKey := imgURL[strings.LastIndex(imgURL, "/")+1:]
	imgKey = strings.TrimSuffix(imgKey, filepath.Ext(imgKey))
	subKey := subURL[strings.LastIndex(subURL, "/")+1:]
	subKey = strings.TrimSuffix(subKey, filepath.Ext(subKey))
	return imgKey, subKey
}
