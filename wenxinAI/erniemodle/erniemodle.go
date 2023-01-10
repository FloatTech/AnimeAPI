// Package erniemodel 百度文心AI大模型
package erniemodel

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/web"
)

type tokendata struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// GetToken 获取当天的token
//
// 申请账号链接:https://wenxin.baidu.com/moduleApi/key
//
// clientID为API key,clientSecret为Secret key
//
// token有效时间为24小时
func GetToken(clientID, clientSecret string) (token string, err error) {
	requestURL := "https://wenxin.baidu.com/moduleApi/portal/api/oauth/token?grant_type=client_credentials&client_id=" + url.QueryEscape(clientID) + "&client_secret=" + url.QueryEscape(clientSecret)
	postData := url.Values{}
	postData.Add("name", "ATRI")
	postData.Add("language", "golang")
	data, err := web.PostData(requestURL, "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}
	var parsed tokendata
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return
	}
	if parsed.Msg != "success" {
		err = errors.New(parsed.Msg + ",code:" + strconv.Itoa(parsed.Code))
		return
	}
	return parsed.Data, nil
}

type parsed struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Result    string `json:"result"`
		RequestID string `json:"requestId"`
	} `json:"data"`
}

// GetResult 创建任务
//
// token:GetToken函数获取,
//
// model:请求类型
//
// txt:用户输入文本
//
// mindeclen:最小生成长度[1, seq_len]
//
// seqlen:最大生成长度[1, 1000]
//
// seqlen决定生成时间:生成512需要16.3s，生成256需要8.1s，生成128需要4.1s
//
// taskprompt:任务类型(非必需)
//
// model：写作文: 1;  写文案: 2;  写摘要: 3;  对对联: 4;  自由问答: 5;  写小说: 6;  补全文本: 7;  自定义: 8;
//
// task_prompt只支持以下：
// PARAGRAPH：引导模型生成一段文章； SENT：引导模型生成一句话； ENTITY：引导模型生成词组； Summarization：摘要； MT：翻译； Text2Annotation：抽取； Correction：纠错； QA_MRC：阅读理解； Dialogue：对话； QA_Closed_book: 闭卷问答； QA_Multi_Choice：多选问答； QuestionGeneration：问题生成； Paraphrasing：复述； NLI：文本蕴含识别； SemanticMatching：匹配； Text2SQL：文本描述转SQL；TextClassification：文本分类； SentimentClassification：情感分析； zuowen：写作文； adtext：写文案； couplet：对对联； novel：写小说； cloze：文
func GetResult(token string, model int, txt string, mindeclen, seqlen int, taskprompt ...string) (result string, err error) {
	requestURL := "https://wenxin.baidu.com/moduleApi/portal/api/rest/1.0/ernie/3.0.2" + strconv.Itoa(model) + "/zeus?" +
		"access_token=" + url.QueryEscape(token)
	postData := url.Values{}
	postData.Add("text", txt)
	postData.Add("min_dec_len", strconv.Itoa(mindeclen))
	postData.Add("seq_len", strconv.Itoa(seqlen))
	postData.Add("topp", "1.0")
	postData.Add("task_prompt", taskprompt[0])
	data, err := web.PostData(requestURL, "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}
	var parsed parsed
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return
	}
	if parsed.Msg != "success" {
		err = errors.New(parsed.Msg + ",code:" + strconv.Itoa(parsed.Code))
		return
	}
	return parsed.Data.Result, nil
}
