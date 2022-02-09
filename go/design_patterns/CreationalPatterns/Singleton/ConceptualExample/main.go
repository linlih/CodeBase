package main

import "fmt"

func main() {
	for i := 0; i < 30; i++ {
		go getInstance()
	}

	// Scanln 和 Scan 类似，但是会在最后一个输入元素为 newline 或者 EOF 结束
	fmt.Scanln()
}
