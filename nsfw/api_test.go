package nsfw

import "testing"

func TestClassify(t *testing.T) {
	p, err := Classify("http://sayuri.fumiama.top/img/?path=詴櫈廘萷泀")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	// t.Fail()
}
