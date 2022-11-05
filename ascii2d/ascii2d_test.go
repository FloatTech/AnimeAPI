package ascii2d

import (
	"crypto/md5"
	"encoding/json"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	r, err := Ascii2d("https://gchat.qpic.cn/gchatpic_new//--05F47960F2546E874F515A403FD174DF/0?term=3")
	if err != nil {
		t.Fatal(err)
	}
	res, err := os.ReadFile("result.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Fatal(err)
	}
	m1 := md5.Sum(res)
	m2 := md5.Sum(b)
	if m1 != m2 {
		f, err := os.Create("result.txt")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = f.Write(b)
		t.Fatal("new result has been written to result.txt")
	}
}
