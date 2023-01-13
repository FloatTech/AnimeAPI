package pixiv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCat(t *testing.T) {
	g, err := Cat(81918463)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, g.Multiple)
	assert.Equal(t, []string{"https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p0.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p1.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p2.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p3.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p4.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p5.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p6.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p7.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p8.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p9.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p10.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p11.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p12.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p13.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p14.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p15.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p16.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p17.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p18.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p19.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p20.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p21.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p22.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p23.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p24.png", "https://i.pximg.net/img-original/img/2020/05/28/19/39/02/81918463_p25.png"}, g.OriginalUrls)
	g, err = Cat(79673672)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, g.Multiple)
	assert.Equal(t, "https://i.pximg.net/img-original/img/2020/02/23/13/26/44/79673672_p0.jpg", g.OriginalURL)
}
