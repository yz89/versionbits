package blockchain

import (
	"fmt"
	"versionbits/chainhash"
)

type blockNode struct {
	hash      chainhash.Hash
	parent    *blockNode
	height    int32
	version   int32
	nonce     uint32
	timestamp int64
}

func HelloBlock() {
	fmt.Println("This is block")
}
