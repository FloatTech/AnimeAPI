package emozi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FloatTech/floatbox/binary"
)

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
