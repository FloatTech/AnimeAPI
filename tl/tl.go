package tl

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

func Translate(target string) (string, error) {
	data, err := web.GetData("http://api.cloolc.club/fanyi?data=" + target)
	if err != nil {
		return "", err
	}
	return binary.BytesToString(binary.NewWriterF(func(w *binary.Writer) {
		for _, v := range gjson.ParseBytes(data).Get("data.0").Get("value").Array() {
			w.WriteString(v.String())
			w.WriteString(",")
		}
	})), nil
}
