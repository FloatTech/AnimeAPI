package main

import (
	"fmt"

	"github.com/FloatTech/AnimeAPI/pixiv"
	"github.com/FloatTech/AnimeAPI/saucenao"
)

func main() {
	id := "75841587"
	data, err := pixiv.Works(id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", data)

	sau, err := saucenao.SauceNAO("http://gchat.qpic.cn/gchatpic_new/213864964/780718903-2150696581-2D6393542C1E07DA915BFEF89ECA0BBD/0?term=2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*sau)
	}
}
