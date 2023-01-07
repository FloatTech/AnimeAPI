package setu

import (
	"os"
	"testing"
	"time"

	"github.com/FloatTech/floatbox/web"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	p, err := NewPool("pool",
		func(s string) (string, error) {
			return "https://pic.moehu.org/large/ec43126fgy1grkj3zrrxsj24dk2k81l1.jpg", nil
		}, web.GetData, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("pool")
	_, err = p.Roll("a")
	if err == nil {
		t.Fatal("unexpected success")
	}
	err = os.Mkdir("pool/a", 0755)
	if err != nil {
		t.Fatal(err)
	}
	s, err := p.Roll("a")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "pool/a/燃相縸沌慀.jpeg", s)
}
