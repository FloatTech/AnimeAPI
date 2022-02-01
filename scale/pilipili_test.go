package scale

import (
	"os"
	"testing"
)

func TestPost(t *testing.T) {
	f, err := os.Open("test/in.png")
	if err != nil {
		t.Fatal(err)
	}
	data, err := Post(f, 4, 4, 2)
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
