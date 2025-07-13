package niu

import "testing"

func TestCreateUserInfoByProps(t *testing.T) {
	user := &userInfo{
		UserID: 123,
		Length: 12,
		WeiGe:  2,
	}
	err := user.applyProp("媚药")
	if err != nil {
		t.Error(err)
	}
	t.Log("成功-----", user)
}

func TestCheckProp(t *testing.T) {
	user := &userInfo{
		UserID: 123,
		Length: 12,
		WeiGe:  2,
	}
	err := user.checkProps("击剑", "jj")
	if err != nil {
		t.Error(err)
	}
	t.Log("成功")
}
