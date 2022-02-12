package baidutts

import (
	"github.com/FloatTech/zbputils/process"
	"os"
)

func init() {
	go func() {
		process.SleepAbout1sTo2s()
		_ = os.MkdirAll(dbpath, 0755)
		os.RemoveAll(cachePath)
		_ = os.MkdirAll(cachePath, 0755)
	}()
}
