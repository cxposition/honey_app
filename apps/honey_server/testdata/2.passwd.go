package main

import (
	"fmt"
	"honey_app/apps/honey_server/utils/pwd"
)

func main() {
	hashPwd, _ := pwd.GenerateFromPassword("1234")
	fmt.Println(hashPwd)
	fmt.Println(pwd.CompareHashAndPassword(hashPwd, "1234"))
}
