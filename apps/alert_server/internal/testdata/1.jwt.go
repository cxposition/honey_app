package main

import (
	"fmt"
	"honey_server/internal/core"
	"honey_server/internal/global"
	"honey_server/internal/utils/jwts"
)

func main() {
	global.Config = core.ReadConfig()
	token, _ := jwts.GetToken(jwts.ClaimsUserInfo{
		UserID: 1,
		Role:   1,
	})
	fmt.Printf("token:%+v\n", token)
	claims, _ := jwts.ParseToken(token)
	fmt.Printf("claims:%+v\n", claims)
}
