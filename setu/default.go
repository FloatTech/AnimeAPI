package setu

import (
	"os"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
)

const DefaultPoolDir = "data/setupool"

var DefaultPool = func() *Pool {
	p, err := NewPool(DefaultPoolDir, func(s string) (string, error) {
		typ := DefaultPoolDir + "/" + s
		if file.IsNotExist(typ) {
			err := os.MkdirAll(typ, 0755)
			if err != nil {
				return "", err
			}
		}
		return "https://img.moehu.org/pic.php?id=pc", nil
	}, web.GetData, time.Minute)
	if err != nil {
		panic(err)
	}
	return p
}()
