package baidutts

import (
	"fmt"
	"os"
	"testing"
)

func TestGetWav(t *testing.T) {
	_ = os.MkdirAll(dbpath, 0755)
	os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	per := 4
	uid := int64(123456)
	tex := "你好，我是超威蓝猫"
	tok := getToken()
	filename := getWav(tex, tok, 5, per, 5, 5, uid)
	fmt.Println(filename)
}
