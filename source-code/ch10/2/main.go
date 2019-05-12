package main

import (
	"fmt"
)

func main() {
	s := "Go语言四十二章经"
	for k, v := range s {
		fmt.Printf("k：%d,v：%c == %d\n", k, v, v)
	}
}
