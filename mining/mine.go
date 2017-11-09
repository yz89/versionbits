package mining

import (
	"fmt"
	"time"
	"versionbits/blockchain"
)

const (
	maxNonce = ^uint32(0) // 2^32 - 1
)

func proofOfWork(header *blockchain.BlockHeader) (bool, uint32) {
	targetDifficulty := blockchain.CompactToBig(header.Bits)
	nonce := uint32(0)
	for ; nonce < maxNonce; nonce++ {
		header.Nonce = nonce
		hash := header.HashBlock()
		bigIntHash := blockchain.HashToBig(&hash)
		// fmt.Println(bigIntHash, targetDifficulty, nonce)
		if bigIntHash.Cmp(targetDifficulty) <= 0 {
			fmt.Printf("%s ", hash)
			return true, nonce
		}
	}
	return false, 0
}

// Mine Yes, let's mining
func Mine() {
	fmt.Println("Mining ...")
	genesisBlock := blockchain.GetGenesisBlock()
	fmt.Println("Genesis Block: ", genesisBlock)

	prevBlock := genesisBlock

	for {
		bits := blockchain.CalcNextRequiredDifficulty(prevBlock)

		nextBlock := prevBlock.GenerateNextBlock()
		nextBlock.Bits = bits

		startTime := time.Now()
		nextBlockHeader := nextBlock.Header()
		solved, nonce := proofOfWork(nextBlockHeader)
		endTime := time.Now()
		if solved {
			elapsedTime := endTime.Sub(startTime).Seconds()
			hashPower := float64(nonce) / (elapsedTime * 1000 * 1000)
			fmt.Printf("Height: %d Version: %b Bits: %x ElapsedTime: %.3fs HashPower: %.2f MH Nonce %d \n", nextBlock.Height, nextBlock.Version, bits, elapsedTime, hashPower, nonce)
			prevBlock = nextBlock
		}
	}
}
