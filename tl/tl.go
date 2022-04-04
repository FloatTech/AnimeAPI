package tl

import (
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
)

func Translate(target string) (string, error) {
	data, err := web.GetData("https://api.cloolc.club/fanyi?data=" + target)
	if err != nil {
		return "", err
	}
	return binary.BytesToString(binary.NewWriterF(func(w *binary.Writer) {
		for _, v := range gjson.ParseBytes(data).Get("data.0").Get("value").Array() {
			s := v.String()
			if len(s) == 0 {
				w.WriteString(",")
				continue
			}
			w.WriteString(s)
		}
	})), nil
}
