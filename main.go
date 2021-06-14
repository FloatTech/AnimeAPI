package main

import (
	"fmt"

	"github.com/FloatTech/AnimeAPI/saucenao"
)

func main() {
	// pixiv
	// temp, _ := pixiv.Works(90281866)
	// fmt.Println(string(temp))
	// saucenao
	sau, err := saucenao.SauceNAO("http://gchat.qpic.cn/gchatpic_new/213864964/780718903-2150696581-2D6393542C1E07DA915BFEF89ECA0BBD/0?term=2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*sau)
	}
}
