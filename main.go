package main

import (
	"fmt"
	"versionbits/blockchain"
	"versionbits/chainhash"
	"versionbits/pow"
)

func main() {
	fmt.Println("Hello Versionbits")
	pow.Mining()
	blockchain.HelloBlock()

	s := "123"
	hash := chainhash.HashH([]byte(s))

	fmt.Println(hash.String())
}
