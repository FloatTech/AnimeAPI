package mockingbird

import (
	"github.com/FloatTech/zbputils/file"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/FloatTech/zbputils/process"
)

// 加载数据库
func init() {
	go func() {
		process.SleepAbout1sTo2s()
		_ = os.MkdirAll(dbpath, 0755)
		os.RemoveAll(cachePath)
		_ = os.MkdirAll(cachePath, 0755)
		_, err := file.GetLazyData(Azfile, false, true)
		if err != nil {
			panic(err)
		}
		_, err = file.GetLazyData(Ysgfile, false, true)
		if err != nil {
			panic(err)
		}
		logrus.Infoln("[mockingbird]加载实例音频")
	}()
}
