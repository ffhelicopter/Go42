package main

import "fmt"

func get() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) // 显示: 10000 10000 数组首字节地址
	return raw[:3]                           // 10000个字节实际只需要引用3个，其他空间浪费
}

func main() {
	data := get()
	fmt.Println(len(data), cap(data), &data[0]) // 显示: 3 10000 数组首字节地址
}
