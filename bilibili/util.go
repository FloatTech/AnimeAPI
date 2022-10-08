package bilibili

import (
	"net/http"
	"strconv"
)

// HumanNum 格式化人数
func HumanNum(res int) string {
	if res/10000 != 0 {
		return strconv.FormatFloat(float64(res)/10000, 'f', 2, 64) + "万"
	}
	return strconv.Itoa(res)
}

// GetRealUrl 获取跳转后的链接
func GetRealUrl(url string) (realurl string, err error) {
	data, err := http.Head(url)
	if err != nil {
		return
	}
	_ = data.Body.Close()
	realurl = data.Request.URL.String()
	return
}
