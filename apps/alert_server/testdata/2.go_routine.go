package main

import (
	"fmt"
	"sync"
)

var lock sync.Mutex

// 字节/阿里/腾讯内部统一写法（2024-2025 版）

func run() {
	panic("panic hhh")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获到panic：", err)
		}
	}()

	run()

	fmt.Println("program over!")
}
