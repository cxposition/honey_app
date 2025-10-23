package main

import (
	"fmt"
	"honey_node/internal/utils/info"
)

func main() {
	list, err := info.GetNetworkList("br")
	fmt.Println(list, err)
}
