package emozi

import "testing"

func TestAll(t *testing.T) {
	usr := Anonymous()
	in := "你好，世界！"
	out, _, err := usr.Marshal(false, in)
	if err != nil {
		t.Fatal(err)
	}
	exp := "🥛‎👔⁡🐴‌👤🌹🐱🐴👩，💦🌞😨🌍➕👴😨👨‍🌾！" //nolint: go-staticcheck
	if out != exp {
		t.Fatal("expected", exp, "but got", out)
	}
	out, err = usr.Unmarshal(false, out)
	if err != nil {
		t.Fatal(err)
	}
	exp = "[你|儗]好，世[界|畍]！"
	if out != exp {
		t.Fatal("expected", exp, "but got", out)
	}
}
