package blockchain

import (
	"encoding/binary"
	"versionbits/chainhash"
)

type BlockHeader struct {
	Version    uint32
	PrevBlock  chainhash.Hash
	MerkleRoot chainhash.Hash
	Timestamp  int64
	Bits       uint32
	Nonce      uint32
}

func (b BlockHeader) serialize() []byte {
	var buf [224]byte
	binary.LittleEndian.PutUint32(buf[0:32], uint32(b.Version))
	binary.LittleEndian.PutUint32(buf[32:64], binary.LittleEndian.Uint32(b.PrevBlock.CloneBytes()))
	binary.LittleEndian.PutUint32(buf[64:96], binary.LittleEndian.Uint32(b.MerkleRoot.CloneBytes()))
	binary.LittleEndian.PutUint64(buf[96:160], uint64(b.Timestamp))
	binary.LittleEndian.PutUint32(buf[160:192], uint32(b.Bits))
	binary.LittleEndian.PutUint32(buf[192:224], uint32(b.Nonce))
	return buf[:]
}

func (b BlockHeader) HashBlock() chainhash.Hash {
	blockBuf := b.serialize()
	hash := chainhash.HashH(blockBuf)
	return hash
}
