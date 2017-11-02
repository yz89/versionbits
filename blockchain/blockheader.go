package blockchain

import (
	"encoding/binary"
	"versionbits/chainhash"
)

type BlockHeader struct {
	Version   uint32
	PrevBlock chainhash.Hash
	Timestamp int64
	Bits      uint32
	Nonce     uint32
}

func (b BlockHeader) serialize() []byte {
	var buf [192]byte
	binary.LittleEndian.PutUint32(buf[0:32], uint32(b.Version))
	binary.LittleEndian.PutUint32(buf[32:64], binary.LittleEndian.Uint32(b.PrevBlock.CloneBytes()))
	binary.LittleEndian.PutUint64(buf[64:128], uint64(b.Timestamp))
	binary.LittleEndian.PutUint32(buf[128:160], uint32(b.Bits))
	binary.LittleEndian.PutUint32(buf[160:192], uint32(b.Nonce))
	return buf[:]
}

func (b BlockHeader) HashBlock() chainhash.Hash {
	blockBuf := b.serialize()
	hash := chainhash.HashH(blockBuf)
	return hash
}
