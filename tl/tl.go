// Package tl 翻译api
package tl

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

// Translate ...
func Translate(target string) (string, error) {
	data, err := web.GetData("http://api.cloolc.club/fanyi?data=" + target)
	if err != nil {
		return "", err
	}
	return binary.BytesToString(binary.NewWriterF(func(w *binary.Writer) {
		meanings := gjson.ParseBytes(data).Get("data.0").Get("value").Array()
		if len(meanings) == 0 {
			w.WriteString("ERROR: 无返回")
			return
		}
		w.WriteString(meanings[0].String())
		for _, v := range meanings[1:] {
			w.WriteString(", ")
			w.WriteString(v.String())
		}
	})), nil
}
