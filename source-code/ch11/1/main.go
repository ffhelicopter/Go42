package main

import (
	"fmt"
)

func main() {

	var arr1 = new([5]int)
	arr := arr1
	arr1[2] = 100
	fmt.Println(arr1[2], arr[2])

	var arr2 [5]int
	newarr := arr2
	arr2[2] = 100
	fmt.Println(arr2[2], newarr[2])
}
