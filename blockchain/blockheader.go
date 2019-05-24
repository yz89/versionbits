package blockchain

import (
	"encoding/binary"
	"time"
	"versionbits/chainhash"
)

type BlockHeader struct {
	Version    uint32
	PrevBlock  chainhash.Hash
	MerkleRoot chainhash.Hash
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

func (b BlockHeader) serialize() []byte {
	var buf [80]byte
	binary.LittleEndian.PutUint32(buf[0:4], uint32(b.Version))
	copy(buf[4:36], b.PrevBlock.CloneBytes())
	copy(buf[36:68], b.MerkleRoot.CloneBytes())
	binary.LittleEndian.PutUint32(buf[68:72], uint32(b.Timestamp.Unix()))
	binary.LittleEndian.PutUint32(buf[72:76], uint32(b.Bits))
	binary.LittleEndian.PutUint32(buf[76:80], uint32(b.Nonce))
	return buf[:]
}

func (b BlockHeader) HashBlock() chainhash.Hash {
	blockBuf := b.serialize()
	hash := chainhash.HashH(blockBuf)
	return hash
}
