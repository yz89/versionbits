package mining

import (
	"fmt"
	"sync"
	"time"
	"versionbits/blockchain"
)

const (
	maxNonce = ^uint32(0) // 2^32 - 1
)

var (
	chSolved   = make(chan *blockchain.BlockNode)
	chNewBlock = make(chan uint32)
	chain      Chain
)

type Miner struct {
	ID               uint32
	chNewBlockHeight chan uint32
}

type Chain struct {
	stateLock sync.RWMutex
	tipBlock  *blockchain.BlockNode
}

func (c *Chain) tip() *blockchain.BlockNode {
	c.stateLock.Lock()
	tipBlock := c.tipBlock
	c.stateLock.Unlock()

	return tipBlock
}

func (c *Chain) addBest(bestBlock *blockchain.BlockNode) bool {
	res := false
	c.stateLock.Lock()
	if bestBlock.Parent == nil || (bestBlock.Parent.Hash == c.tipBlock.Hash && bestBlock.Height == c.tipBlock.Height+1) {
		c.tipBlock = bestBlock
		res = true
	}
	c.stateLock.Unlock()
	return res
}

func solveBlock(node *blockchain.BlockNode) bool {
	header := node.Header()
	targetDifficulty := blockchain.CompactToBig(header.Bits)
	nonce := uint32(0)
	for ; nonce < maxNonce; nonce++ {
		bestBlock := chain.tip()
		if bestBlock.Hash != node.Parent.Hash {
			// 当前工作不是在最长链上挖，停止，进入下一轮
			return false
		}
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

func mine() {
	for {
		bestBlock := chain.tip()
		nextBlock := bestBlock.GenerateNextBlock()
		startTime := time.Now()
		solvedBlock := solveBlock(nextBlock)
		endTime := time.Now()

		if solvedBlock {
			chain.addBest(nextBlock)
			elapsedTime := endTime.Sub(startTime).Seconds()
			hashPower := float64(nextBlock.Nonce) / (elapsedTime * 1000 * 1000)
			fmt.Printf("%s Height: %d Version: %b Bits: %x ElapsedTime: %.3fs HashPower: %.2f MH Nonce %d \n", nextBlock.Hash, nextBlock.Height, nextBlock.Version, nextBlock.Bits, elapsedTime, hashPower, nextBlock.Nonce)
		}
	}
}

// Start let's mining
func Start() {
	fmt.Println("Start Mining ...")
	genesisBlock := blockchain.GetGenesisBlock()
	fmt.Println("Genesis Block: ", genesisBlock.Height, genesisBlock.Version, genesisBlock.Bits, genesisBlock.Nonce)

	chain.addBest(genesisBlock)

	mine()
}
