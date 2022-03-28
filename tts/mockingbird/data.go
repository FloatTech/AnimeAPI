package mockingbird

import (
	"os"
)

// 加载数据库
func init() {
	_ = os.MkdirAll(dbpath, 0755)
	os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
}
