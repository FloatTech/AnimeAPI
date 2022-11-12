package qzone

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/FloatTech/floatbox/binary"
)

var (
	cookie = "RK=u8/BwhugZ1; ptcz=a8641719c2bac979d89df84627a14407d3947821484c23765b1f692443c0d23f; pgv_pvid=6474350135; iip=0; logTrackKey=b4c60d611df64dfaa180ebd43aea3c6f; __Q_w_s_hat_seed=1; _tc_unionid=e66b4f13-ea92-4afb-b5d0-6694f8c2b443; pac_uid=1_1156544355; QZ_FE_WEBP_SUPPORT=1; __Q_w_s__QZN_TodoMsgCnt=1; o_cookie=1156544355; feeds_selector=2; zzpaneluin=; zzpanelkey=; _qpsvr_localtk=0.33152714365333447; pgv_info=ssid=s1907750736; uin=o1776620359; skey=@vwaG3O59W; p_uin=o1776620359; pt4_token=lRUQWfNSsWVim2G5hdqtVo2rtuGIyBZS26Bjf08gcJ8_; p_skey=iywhVczQpsvg6hcUH8bty53Ge0SBcdkBN7nM5eBp-Z0_; Loading=Yes; cpu_performance_v8=15"
)

func TestManager_PublishEmotion(t *testing.T) {
	type args struct {
		Content string
	}
	m := NewManager(cookie)
	gotResult, err := m.EmotionPublish("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}

func TestManager_UploadImage(t *testing.T) {
	m := NewManager(cookie)
	path := `D:\Documents\Pictures\日南葵\日南葵1.jpg`
	srcByte, err := os.ReadFile(path)
	if err != nil {
		return
	}
	picBase64 := base64.StdEncoding.EncodeToString(srcByte)
	if err != nil {
		t.Fatal(err)
	}
	gotResult, err := m.UploadImage(picBase64)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}

func TestManager_Msglist(t *testing.T) {
	m := NewManager(cookie)
	gotResult, err := m.EmotionMsglist("20", "100")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}

func TestLogin(t *testing.T) {
	var (
		qrsig           string
		ptqrtoken       string
		ptqrloginCookie string
		redirectCookie  string
		data            []byte
		err             error
	)
	data, qrsig, ptqrtoken, err = Ptqrshow()
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("ptqrcode.png", data, 0666)
	if err != nil {
		t.Fatal(err)
	}
LOOP:
	for {
		time.Sleep(2 * time.Second)
		data, ptqrloginCookie, err = Ptqrlogin(qrsig, ptqrtoken)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("ptqrloginCookie:%v\n", ptqrloginCookie)
		text := binary.BytesToString(data)
		t.Logf("text:%v\n", text)
		switch {
		case strings.Contains(text, "二维码已失效"):
			t.Fatal("二维码已失效, 登录失败")
			return
		case strings.Contains(text, "登录成功"):
			_ = os.Remove("ptqrcode.png")
			dealedCheckText := strings.ReplaceAll(text, "'", "")
			redirectURL := strings.Split(dealedCheckText, ",")[2]
			redirectCookie, err = LoginRedirect(redirectURL)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("ptqrloginCookie:%v\n", redirectCookie)
			break LOOP
		}
	}
	m := NewManager(ptqrloginCookie + redirectCookie)
	t.Logf("m:%#v\n", m)
	err = os.WriteFile("cookie.txt", binary.StringToBytes(ptqrloginCookie+redirectCookie), 0666)
	if err != nil {
		t.Fatal(err)
	}
	gotResult, err := m.EmotionPublish("真好", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gotResult:%#v\n", gotResult)
}
