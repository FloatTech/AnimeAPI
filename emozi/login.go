package emozi

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/tidwall/gjson"
)

// NewUser 创建已注册用户实例
//
// 注册请前往 API 网址
func NewUser(name, pswd string) (usr User) {
	usr.name = name
	usr.pswd = pswd
	return
}

// Anonymous 创建匿名用户
//
// 有访问请求数限制
func Anonymous() (usr User) {
	return
}

// Login 登录
func (usr *User) Login() error {
	data, err := web.GetData(api + "getLoginSalt?username=" + url.QueryEscape(usr.name))
	if err != nil {
		return err
	}
	r := gjson.ParseBytes(data)
	if r.Get("code").Int() != 0 {
		return errors.New(r.Get("message").Str)
	}
	salt := r.Get("result.salt").Str
	h := md5.New()
	h.Write([]byte(usr.pswd))
	h.Write([]byte(salt))
	passchlg := hex.EncodeToString(h.Sum(make([]byte, 0, md5.Size)))
	w := binary.SelectWriter()
	defer binary.PutWriter(w)
	err = json.NewEncoder(w).Encode(&loginbody{
		Username: usr.name,
		Password: passchlg,
		Salt:     salt,
	})
	if err != nil {
		return err
	}
	data, err = web.PostData(api+"login", "application/json", (*bytes.Buffer)(w))
	if err != nil {
		return err
	}
	r = gjson.ParseBytes(data)
	if r.Get("code").Int() != 0 {
		return errors.New(r.Get("message").Str)
	}
	usr.auth = r.Get("result.token").Str
	return nil
}
