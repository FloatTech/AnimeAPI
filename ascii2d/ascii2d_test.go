package ascii2d

import (
	"crypto/md5"
	"encoding/json"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	r, err := Ascii2d("https://gchat.qpic.cn/gchatpic_new//663582230-2596768878-CECBADF39E266F89655249A56810EA4F/0")
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
		f.Write(b)
		t.Fatal("new result has been written to result.txt")
	}
}
