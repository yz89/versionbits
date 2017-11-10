package mining

import (
	"fmt"
	"time"
	"versionbits/blockchain"
)

const (
	maxNonce = ^uint32(0) // 2^32 - 1
)

var (
	chSolved   = make(chan *blockchain.BlockNode)
	chNewBlock = make(chan uint32)
)

type Miner struct {
	ID               uint32
	chNewBlockHeight chan uint32
}

func solveBlock(node *blockchain.BlockNode) bool {
	header := node.Header()

	targetDifficulty := blockchain.CompactToBig(header.Bits)
	nonce := uint32(0)
	for ; nonce < maxNonce; nonce++ {
		header.Nonce = nonce
		hash := header.HashBlock()
		bigIntHash := blockchain.HashToBig(&hash)
		if bigIntHash.Cmp(targetDifficulty) <= 0 {
			// fmt.Printf("%d ", m.ID)
			node.Hash = hash
			node.Nonce = nonce
			return true
		}
	}
	return false
}

// Start let's mining
func Start() {
	fmt.Println("Start Mining ...")
	genesisBlock := blockchain.GetGenesisBlock()
	fmt.Println("Genesis Block: ", genesisBlock.Height, genesisBlock.Version, genesisBlock.Bits, genesisBlock.Nonce)

	prevBlock := genesisBlock

	for {
		nextBlock := prevBlock.GenerateNextBlock()
		startTime := time.Now()
		solvedBlock := solveBlock(nextBlock)
		endTime := time.Now()

		if solvedBlock {
			elapsedTime := endTime.Sub(startTime).Seconds()
			hashPower := float64(nextBlock.Nonce) / (elapsedTime * 1000 * 1000)
			fmt.Printf("%s Height: %d Version: %b Bits: %x ElapsedTime: %.3fs HashPower: %.2f MH Nonce %d \n", nextBlock.Hash, nextBlock.Height, nextBlock.Version, nextBlock.Bits, elapsedTime, hashPower, nextBlock.Nonce)
			prevBlock = nextBlock
		}
	}
}
