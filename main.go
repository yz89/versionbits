package main

import (
	"fmt"
	"versionbits/blockchain"
	"versionbits/pow"
)

func main() {
	fmt.Println("Hello Versionbits")
	pow.Mining()
	blockchain.HelloBlock()
}
