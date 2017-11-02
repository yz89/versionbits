package blockchain

import (
	"encoding/binary"
	"fmt"
	"time"
	"versionbits/chainhash"
)

type BlockHeader struct {
	PrevBlock chainhash.Hash
	Version   uint32
	Nonce     uint32
	Timestamp int64
}

type BlockNode struct {
	parent    *BlockNode
	hash      chainhash.Hash
	Height    int32
	Version   int32
	Nonce     uint32
	Timestamp int64
}

func (b BlockHeader) serialize() []byte {
	var buf [160]byte
	binary.LittleEndian.PutUint32(buf[0:32], binary.LittleEndian.Uint32(b.PrevBlock.CloneBytes()))
	binary.LittleEndian.PutUint32(buf[32:64], uint32(b.Version))
	binary.LittleEndian.PutUint32(buf[64:96], uint32(b.Nonce))
	binary.LittleEndian.PutUint64(buf[96:160], uint64(b.Timestamp))
	return buf[:]
}

func (b BlockHeader) HashBlock() chainhash.Hash {
	blockBuf := b.serialize()
	hash := chainhash.HashH(blockBuf)
	return hash
}

func (b BlockHeader) GenerateNextBlock() BlockHeader {
	var newBlock = BlockHeader{
		PrevBlock: b.HashBlock(),
		Version:   0,
		Nonce:     0,
		Timestamp: time.Now().Unix(),
	}

	return newBlock
}

func (b BlockHeader) SetNonce(nonce uint32) BlockHeader {
	b.Nonce = nonce
	return b
}

func GetGenesisBlock() BlockHeader {
	zeroHash := chainhash.HashH([]byte{0})
	var genesisBlock = BlockHeader{
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
