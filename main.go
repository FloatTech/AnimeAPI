package main

import (
	"fmt"

	"github.com/FloatTech/AnimeAPI/pixiv"
)

func main() {
	temp, _ := pixiv.Works(90281866)
	fmt.Println(string(temp))
}
