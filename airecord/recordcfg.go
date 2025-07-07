package airecord

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	recordcfg  recordconfig
	configPath = "data/airecord/recordconfig.json" // 配置文件路径
)

func init() {
	if err := loadConfig(); err != nil {
		logrus.Warnln("[airecord] WARN: 加载配置文件失败，使用默认配置:", err)
	} else {
		logrus.Infoln("[airecord] 成功从文件加载语音记录配置")
	}
}

// recordconfig 存储语音记录相关配置
type recordconfig struct {
	ModelName string `json:"modelName"` // 语音模型名称
	ModelID   string `json:"modelID"`   // 语音模型ID
	Customgid int64  `json:"customgid"` // 自定义群ID
}

// GetRecordConfig 返回当前语音记录配置信息
func GetRecordConfig() recordconfig {
	return recordcfg
}

// SetRecordModel 设置语音记录模型
func SetRecordModel(modelName, modelID string) {
	recordcfg.ModelName = modelName
	recordcfg.ModelID = modelID
	saveConfig() // 保存配置
}

// SetCustomGID 设置自定义群ID
func SetCustomGID(gid int64) {
	recordcfg.Customgid = gid
	saveConfig() // 保存配置
}

// PrintRecordConfig 生成格式化的语音记录配置信息字符串
func PrintRecordConfig(recCfg recordconfig) string {
	var builder strings.Builder
	builder.WriteString("当前语音记录配置：\n")
	builder.WriteString(fmt.Sprintf("• 语音模型名称：%s\n", recCfg.ModelName))
	builder.WriteString(fmt.Sprintf("• 语音模型ID：%s\n", recCfg.ModelID))
	builder.WriteString(fmt.Sprintf("• 自定义群ID：%d\n", recCfg.Customgid))
	return builder.String()
}

// saveConfig 将配置保存到JSON文件
func saveConfig() error {
	data, err := json.MarshalIndent(recordcfg, "", "  ")
	if err != nil {
		logrus.Warnln("ERROR: 序列化配置失败:", err)
		return err
	}
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		logrus.Warnln("ERROR: 写入配置文件失败:", err)
		return err
	}
	return nil
}

// loadConfig 从JSON文件加载配置
func loadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &recordcfg)
	if err != nil {
		logrus.Warnln("ERROR: 解析配置文件失败:", err)
		return err
	}
	return nil
}
