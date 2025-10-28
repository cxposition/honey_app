package main

import (
	"fmt"
	"honey_node/internal/utils/info"
)

func main() {
	list, err := info.GetNetworkList([]string{"br-"})
	fmt.Println(list, err)
}
