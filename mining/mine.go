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
	wg         sync.WaitGroup
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
	ok := false
	c.stateLock.Lock()
	if bestBlock.Parent == nil || (bestBlock.Parent.Hash == c.tipBlock.Hash && bestBlock.Height == c.tipBlock.Height+1) {
		c.tipBlock = bestBlock
		ok = true
	}
	c.stateLock.Unlock()
	return ok
}

func solveBlock(node *blockchain.BlockNode, ticker *time.Ticker) bool {
	header := node.Header()
	targetDifficulty := blockchain.CompactToBig(header.Bits)
	nonce := uint32(0)
	for ; nonce < maxNonce; nonce++ {

		// 定时检查当前工作是否在最长链上挖
		select {
		case <-ticker.C:
			bestBlock := chain.tip()
			if bestBlock.Hash != node.Parent.Hash {
				// 没有在最长链上挖，停止，重新获取最新块
				return false
			}
		default:
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

func mine(minerID uint32) {
	for {
		bestBlock := chain.tip()
		nextBlock := bestBlock.GenerateNextBlock()
		ticker := time.NewTicker(50 * time.Millisecond)

		startTime := time.Now()
		solvedBlock := solveBlock(nextBlock, ticker)
		endTime := time.Now()

		if solvedBlock {
			ok := chain.addBest(nextBlock)
			if ok {
				elapsedTime := endTime.Sub(startTime).Seconds()
				hashPower := float64(nextBlock.Nonce) / (elapsedTime * 1000 * 1000)
				fmt.Printf("%d %s Height: %d Version: %b Bits: %x ElapsedTime: %.3fs HashPower: %.2f MH Nonce %d \n", minerID, nextBlock.Hash, nextBlock.Height, nextBlock.Version, nextBlock.Bits, elapsedTime, hashPower, nextBlock.Nonce)
			} else {
				fmt.Printf("%d %s Height: %d Abandon!! \n", minerID, nextBlock.Hash, nextBlock.Height)
			}
		}
	}
}

// Start let's mining
func Start() {
	fmt.Println("Start Mining ...")
	genesisBlock := blockchain.GetGenesisBlock()
	fmt.Println("Genesis Block: ", genesisBlock.Height, genesisBlock.Version, genesisBlock.Bits, genesisBlock.Nonce)

	chain.addBest(genesisBlock)

	wg.Add(1)

	go mine(1)
	go mine(2)
	go mine(3)
	go mine(4)

	wg.Wait()
}
