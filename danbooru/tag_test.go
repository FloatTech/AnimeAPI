package danbooru

import (
	"testing"
)

func TestRender(t *testing.T) {
	u := "http://192.168.7.81:62002/img?arg=get&name=卆悥什繐捀.webp"
	c, err := TagURL("俁罰穬歭俀", u)
	if err != nil {
		t.Fatal(err)
	}
	c.Canvas.SavePNG("out.png")
}
