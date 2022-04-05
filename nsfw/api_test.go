package nsfw

import "testing"

func TestClassify(t *testing.T) {
	p, err := Classify("https://1mg.obfs.dev/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	// t.Fail()
}
