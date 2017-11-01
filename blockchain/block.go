package blockchain

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
	"versionbits/chainhash"
)

type blockHeader struct {
	PrevBlock chainhash.Hash
	Version   uint32
	Nonce     uint32
	Timestamp int64
}

type blockNode struct {
	parent    *blockNode
	hash      chainhash.Hash
	Height    int32
	Version   int32
	Nonce     uint32
	Timestamp int64
}

func (b blockHeader) serialize() []byte {
	var buf [160]byte
	binary.LittleEndian.PutUint32(buf[0:32], binary.LittleEndian.Uint32(b.PrevBlock.CloneBytes()))
	binary.LittleEndian.PutUint32(buf[32:64], uint32(b.Version))
	binary.LittleEndian.PutUint32(buf[64:96], uint32(b.Nonce))
	binary.LittleEndian.PutUint64(buf[96:160], uint64(b.Timestamp))
	return buf[:]
}

func (b blockHeader) HashBlock() chainhash.Hash {
	blockBuf := b.serialize()
	hash := chainhash.HashH(blockBuf)
	return hash
}

func (b blockHeader) GenerateNextBlock() blockHeader {
	var newBlock = blockHeader{
		PrevBlock: b.HashBlock(),
		Version:   0,
		Nonce:     rand.Uint32(),
		Timestamp: time.Now().Unix(),
	}

	return newBlock
}

func GetGenesisBlock() blockHeader {
	zeroHash := chainhash.HashH([]byte{0})
	var genesisBlock = blockHeader{
		PrevBlock: zeroHash,
		Version:   0,
		Nonce:     0,
		Timestamp: time.Now().Unix(),
	}

	return genesisBlock
}

func HelloBlock() {
	fmt.Println("This is block")
}
