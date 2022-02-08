package danbooru

import (
	"testing"

	"github.com/FloatTech/zbputils/img/writer"
)

func TestRender(t *testing.T) {
	u := "http://192.168.7.81:62002/img?arg=get&name=卆悥什繐捀.webp"
	im, err := TagURL("俁罰穬歭俀", u)
	if err != nil {
		t.Fatal(err)
	}
	writer.SavePNG2Path("out.png", im)
}
