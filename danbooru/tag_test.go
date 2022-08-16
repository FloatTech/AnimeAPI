package danbooru

import (
	"testing"

	"github.com/FloatTech/floatbox/img/writer"
)

func TestRender(t *testing.T) {
	u := "https://1mg.obfs.dev/"
	im, err := TagURL("random", u)
	if err != nil {
		t.Fatal(err)
	}
	writer.SavePNG2Path("out.png", im)
}
