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
	chain blockchain.Chain
	wg    sync.WaitGroup
)

type Miner struct {
	ID               uint32
	chNewBlockHeight chan uint32
}

func solveBlock(node *blockchain.BlockNode, ticker *time.Ticker) bool {
	header := node.Header()
	targetDifficulty := blockchain.CompactToBig(header.Bits)
	nonce := uint32(0)
	for ; nonce < maxNonce; nonce++ {

		// 定时检查当前工作是否在最长链上挖
		select {
		case <-ticker.C:
			bestBlock := chain.Tip()
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
			node.Hash = hash
			node.Nonce = nonce
			return true
		}
	}
	return false
}

func mine(minerID uint32) {
	for {
		bestBlock := chain.Tip()
		nextBlock := bestBlock.GenerateNextBlock()
		nextBlock.Version = blockchain.GetNextVersion(minerID)
		ticker := time.NewTicker(100 * time.Millisecond)

		startTime := time.Now()
		solvedBlock := solveBlock(nextBlock, ticker)
		endTime := time.Now()

		if solvedBlock {
			ok := chain.AddBest(nextBlock)
			if ok {
				elapsedTime := endTime.Sub(startTime).Seconds()
				hashPower := float64(nextBlock.Nonce) / (elapsedTime * 1000 * 1000)
				fmt.Printf("%d %s Height: %d Version: %x Bits: %x ElapsedTime: %.3fs HashPower: %.2f MH Nonce %d \n", minerID, nextBlock.Hash, nextBlock.Height, nextBlock.Version, nextBlock.Bits, elapsedTime, hashPower, nextBlock.Nonce)
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

	chain.AddBest(genesisBlock)

	wg.Add(1)

	go mine(blockchain.Miner1)
	go mine(blockchain.Miner2)
	go mine(blockchain.Miner3)
	go mine(blockchain.Miner4)

	wg.Wait()
}
