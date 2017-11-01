package mining

import (
	"fmt"
	"versionbits/blockchain"
)

// Mining yes, let's mining
func Mining() {
	fmt.Println("Mining ...")
	genesisBlock := blockchain.GetGenesisBlock()
	fmt.Println("Genesis Block: ", genesisBlock)

	prevBlock := genesisBlock

	for {
		nextBlock := prevBlock.GenerateNextBlock()
		hash := nextBlock.HashBlock()

		// targetDifficulty := big.NewInt(0)
		// targetDifficulty.SetString("823530421679738103333889295373420769964125341418332322887543236067889410", 0)
		var bits = uint32(0x20000111)
		targetDifficulty := blockchain.CompactToBig(bits)

		currentHash := blockchain.HashToBig(&hash)

		// fmt.Println(currentHash, targetDifficulty)

		if currentHash.Cmp(targetDifficulty) <= 0 {
			fmt.Printf("%s version %d nonce %d \n", hash, nextBlock.Version, nextBlock.Nonce)
			prevBlock = nextBlock
		}
	}

}
