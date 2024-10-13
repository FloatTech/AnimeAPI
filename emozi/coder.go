package emozi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FloatTech/floatbox/binary"
)

// Marshal 编码
//
//   - randomSameMeaning 随机使用近音颜文字
//   - text 中文文本
//   - choices 多音字选择
func (usr *User) Marshal(randomSameMeaning bool, text string, choices ...int) (string, []int, error) {
	w := binary.SelectWriter()
	defer binary.PutWriter(w)
	err := json.NewEncoder(w).Encode(&encodebody{
		Random: randomSameMeaning,
		Text:   text,
		Choice: choices,
	})
	if err != nil {
		return "", nil, err
	}
	req, err := http.NewRequest("POST", api+"encode", (*bytes.Buffer)(w))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if usr.auth != "" {
		req.Header.Set("Authorization", usr.auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	r := encoderesult{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", nil, err
	}
	if r.Code != 0 {
		return "", nil, errors.New(r.Message)
	}
	return r.Result.Text, r.Result.Choice, nil
}

// Unmarshal 解码
//
//   - force 强制解码不是由程序生成的转写
//   - text 颜文字文本
func (usr *User) Unmarshal(force bool, text string) (string, error) {
	w := binary.SelectWriter()
	defer binary.PutWriter(w)
	err := json.NewEncoder(w).Encode(&decodebody{
		Force: force,
		Text:  text,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", api+"decode", (*bytes.Buffer)(w))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if usr.auth != "" {
		req.Header.Set("Authorization", usr.auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	r := decoderesult{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}
	if r.Code != 0 {
		return "", errors.New(r.Message)
	}
	return r.Result, nil
}
