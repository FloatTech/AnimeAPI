package scale

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	data, err := Get("https://gchat.qpic.cn/gchatpic_new//663582230-2596768878-CECBADF39E266F89655249A56810EA4F/0", 1, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("test/getout.webp")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(data)
	f.Close()
}

func TestPost(t *testing.T) {
	f, err := os.Open("test/0.jpg")
	if err != nil {
		t.Fatal(err)
	}
	data, err := Post(f, 1, 2, 4)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	f, err = os.Create("test/out.webp")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(data)
	f.Close()
}
