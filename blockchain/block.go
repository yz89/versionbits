package blockchain

import "fmt"

type Block struct {
	nHeight  int
	nVersion int
}

func HelloBlock() {
	fmt.Println("This is block")
}
