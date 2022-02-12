package baidutts

import (
	"os"
)

func init() {
	_ = os.MkdirAll(dbpath, 0755)
	_ = os.RemoveAll(cachePath)
	err := os.MkdirAll(cachePath, 0755)
	if err != nil {
		panic(err)
	}
}
