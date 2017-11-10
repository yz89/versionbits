package mining

import (
	"fmt"
	"math/rand"
	"time"
	"versionbits/blockchain"
	"versionbits/chainhash"
)

const (
	maxNonce = ^uint32(0) // 2^32 - 1
)

var (
	chSolved   = make(chan *blockchain.BlockNode)
	chNewBlock = make(chan uint32)

	miners []*Miner
)

type Miner struct {
	ID               uint32
	chNewBlockHeight chan uint32
}

func (m Miner) solveBlock(node blockchain.BlockNode) {
	// simulate different transactions collection
	transactions := []byte{byte(rand.Intn(127)), byte(rand.Intn(127)), byte(rand.Intn(127)), byte(rand.Intn(127))}
	node.MerkleRoot = chainhash.HashH(transactions)

	header := node.Header()

	// fmt.Println(m.ID, header.HashBlock(), transactions)

	targetDifficulty := blockchain.CompactToBig(header.Bits)
	nonce := uint32(0)
	for ; nonce < maxNonce; nonce++ {
		select {
		case <-m.chNewBlockHeight:
			return
		default:
			// Non-blocking select to fall through
		}
		header.Nonce = nonce
		hash := header.HashBlock()
		bigIntHash := blockchain.HashToBig(&hash)
		// fmt.Println(bigIntHash, targetDifficulty, nonce)
		if bigIntHash.Cmp(targetDifficulty) <= 0 {
			fmt.Printf("%d ", m.ID)
			node.Hash = hash
			node.Nonce = nonce

			nextRound(node.Height)
			chSolved <- &node
			return
		}
	}
	return
}

func nextRound(height uint32) {
	for _, miner := range miners {
		miner.chNewBlockHeight <- height
	}
}

// Mine Yes, let's mining
func Mine() {
	fmt.Println("Mining ...")
	genesisBlock := blockchain.GetGenesisBlock()
	fmt.Println("Genesis Block: ", genesisBlock.Height, genesisBlock.Version, genesisBlock.Bits, genesisBlock.Nonce)

	prevBlock := genesisBlock

	for {
		bits := blockchain.CalcNextRequiredDifficulty(prevBlock)

		nextBlock := prevBlock.GenerateNextBlock()
		nextBlock.Bits = bits

		miner1 := Miner{ID: 1, chNewBlockHeight: make(chan uint32)}
		miner2 := Miner{ID: 2, chNewBlockHeight: make(chan uint32)}
		miner3 := Miner{ID: 3, chNewBlockHeight: make(chan uint32)}
		miner4 := Miner{ID: 4, chNewBlockHeight: make(chan uint32)}

		miners = append(miners, &miner1, &miner2, &miner3, &miner4)

		go miner1.solveBlock(*nextBlock)
		go miner2.solveBlock(*nextBlock)
		go miner3.solveBlock(*nextBlock)
		go miner4.solveBlock(*nextBlock)

		startTime := time.Now()
		solvedBlock, ok := <-chSolved
		endTime := time.Now()

		if !ok {
			fmt.Println("not find ")
		}

		if solvedBlock != nil {
			miners = make([]*Miner, 0)
			elapsedTime := endTime.Sub(startTime).Seconds()
			hashPower := float64(solvedBlock.Nonce) / (elapsedTime * 1000 * 1000)
			fmt.Printf("%s Height: %d Version: %b Bits: %x ElapsedTime: %.3fs HashPower: %.2f MH Nonce %d \n", solvedBlock.Hash, solvedBlock.Height, solvedBlock.Version, bits, elapsedTime, hashPower, solvedBlock.Nonce)
			prevBlock = solvedBlock
		}
	}
}
