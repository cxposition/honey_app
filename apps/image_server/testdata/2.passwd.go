package main

import (
	"fmt"
	"image_server/internal/utils/pwd"
)

func main() {
	hashPwd, _ := pwd.GenerateFromPassword("1234")
	fmt.Println(hashPwd)
	fmt.Println(pwd.CompareHashAndPassword(hashPwd, "1234"))
}
