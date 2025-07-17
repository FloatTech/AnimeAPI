// Package airecord 群应用：AI声聊配置
package airecord

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

// Storage 语音记录配置存储
type Storage struct {
	sync.RWMutex
	db sql.Sqlite
}

var (
	sdb = &Storage{
		db: sql.New("data/airecord/recordconfig.db"),
	}
)

func init() {
	if err := os.MkdirAll("data/airecord", 0755); err != nil {
		panic(err)
	}
	if err := sdb.db.Open(time.Hour * 24); err != nil {
		panic(err)
	}
	if err := sdb.db.Create("config", &recordconfig{}); err != nil {
		panic(err)
	}
}

// recordconfig 存储语音记录相关配置
type recordconfig struct {
	ID        int64  `db:"id"`        // 主键ID
	ModelName string `db:"modelName"` // 语音模型名称
	ModelID   string `db:"modelID"`   // 语音模型ID
	Customgid int64  `db:"customgid"` // 自定义群ID
}

// GetConfig 获取当前配置
func GetConfig() recordconfig {
	sdb.RLock()
	defer sdb.RUnlock()
	cfg := recordconfig{}
	_ = sdb.db.Find("config", &cfg, "WHERE id = 1")
	return cfg
}

// SetRecordModel 设置语音记录模型
func SetRecordModel(modelName, modelID string) error {
	cfg := GetConfig()
	sdb.Lock()
	defer sdb.Unlock()
	return sdb.db.Insert("config", &recordconfig{
		ID:        1,
		ModelName: modelName,
		ModelID:   modelID,
		Customgid: cfg.Customgid,
	})
}

// SetCustomGID 设置自定义群ID
func SetCustomGID(gid int64) error {
	cfg := GetConfig()
	sdb.Lock()
	defer sdb.Unlock()
	return sdb.db.Insert("config", &recordconfig{
		ID:        1,
		ModelName: cfg.ModelName,
		ModelID:   cfg.ModelID,
		Customgid: gid,
	})
}

// PrintRecordConfig 生成格式化的语音记录配置信息字符串
func PrintRecordConfig() string {
	cfg := GetConfig()
	var builder strings.Builder
	builder.WriteString("当前语音记录配置：\n")
	builder.WriteString(fmt.Sprintf("• 语音模型名称：%s\n", cfg.ModelName))
	builder.WriteString(fmt.Sprintf("• 语音模型ID：%s\n", cfg.ModelID))
	builder.WriteString(fmt.Sprintf("• 自定义群ID：%d\n", cfg.Customgid))
	return builder.String()
}
