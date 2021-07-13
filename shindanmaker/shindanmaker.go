package shindanmaker

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	xpath "github.com/antchfx/htmlquery"
)

var (
	token  = ""
	cookie = ""
)

// Shindanmaker 基于 https://shindanmaker.com 的 API
// id 是的不同页面的 url 里的数字，例如 https://shindanmaker.com/a/162207 里的 162207
// name 是要被测定的人的名字，影响测定结果
func Shindanmaker(id int64, name string) (string, error) {
	url := fmt.Sprintf("https://shindanmaker.com/%d", id)
	// seed 使每一天的结果都不同
	now := time.Now()
	seed := fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
	name = name + seed

	// 刷新 token 和 cookie
	if token == "" || cookie == "" {
		if err := refresh(); err != nil {
			return "", err
		}
	}

	// 组装参数
	client := &http.Client{}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("_token", token)
	_ = writer.WriteField("shindanName", name)
	_ = writer.WriteField("hiddenName", "名無しのR")
	_ = writer.Close()
	// 发送请求
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Cookie", cookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 解析XPATH
	doc, err := xpath.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	// 取出每个返回的结果
	list := xpath.Find(doc, `//*[@id="shindanResult"]`)
	if len(list) == 0 {
		token = ""
		cookie = ""
		return "", errors.New("无法查找到结果，请稍后再试")
	}
	var output = []string{}
	for child := list[0].FirstChild; child != nil; child = child.NextSibling {
		if text := xpath.InnerText(child); text != "" {
			output = append(output, text)
		} else {
			output = append(output, "\n")
		}
	}
	return strings.ReplaceAll(strings.Join(output, ""), seed, ""), nil
}

// refresh 刷新 cookie 和 token
func refresh() error {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://shindanmaker.com/587874", nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// 获取 cookie
	if temp := resp.Header.Values("Set-Cookie"); len(temp) == 0 {
		return errors.New("刷新 cookie 时发生错误")
	} else {
		cookie = temp[len(temp)-1]
	}
	if !strings.Contains(cookie, "_session") {
		return errors.New("刷新 cookie 时发生错误")
	}
	// 获取 token
	defer resp.Body.Close()
	doc, err := xpath.Parse(resp.Body)
	if err != nil {
		return err
	}
	list := xpath.Find(doc, `//*[@id="shindanForm"]/input`)
	if len(list) == 0 {
		return errors.New("刷新 token 时发生错误")
	}
	token = list[0].Attr[2].Val
	return nil
}
