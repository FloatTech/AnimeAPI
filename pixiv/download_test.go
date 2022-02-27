package pixiv

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadToCache(t *testing.T) {
	illust, err := Works(96415148)
	if err != nil {
		t.Fatal(err)
	}
	err = illust.DownloadToCache(0)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(illust.Path(0))
	if err != nil {
		t.Fatal(err)
	}
	m := md5.Sum(data)
	ms := hex.EncodeToString(m[:])
	assert.Equal(t, "bc635c27d278c414ca8347487d005b6d", ms)
}
