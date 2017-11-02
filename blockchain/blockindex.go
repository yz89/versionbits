package blockchain

import (
	"time"
	"versionbits/chainhash"
)

type BlockNode struct {
	parent    *BlockNode
	hash      chainhash.Hash
	Height    uint32
	Version   uint32
	Timestamp int64
	Nonce     uint32
	Bits      uint32
}

func initBlockNode(node *BlockNode, blockHeader *BlockHeader, height uint32) {
	*node = BlockNode{
		hash:      blockHeader.HashBlock(),
		Height:    height,
		Version:   blockHeader.Version,
		Timestamp: blockHeader.Timestamp,
		Nonce:     blockHeader.Nonce,
		Bits:      blockHeader.Bits,
	}
}

func newBlockNode(blockHeader *BlockHeader, height uint32) *BlockNode {
	var node BlockNode
	initBlockNode(&node, blockHeader, height)
	return &node
}

func (node *BlockNode) Header() *BlockHeader {
	prevHash := chainhash.Hash{}
	if node.parent != nil {
		prevHash = node.parent.hash
	}
	return &BlockHeader{
		Version:   node.Version,
		PrevBlock: prevHash,
		Timestamp: node.Timestamp,
		Bits:      node.Bits,
		Nonce:     node.Nonce,
	}
}

func (node *BlockNode) Ancestor(height uint32) *BlockNode {
	if height < 0 || height > node.Height {
		return nil
	}

	n := node
	for ; n != nil && n.Height != height; n = n.parent {
		// Intentionally left blank
	}
	return n
}

func (node *BlockNode) RelativeAncestor(distance uint32) *BlockNode {
	return node.Ancestor(node.Height - distance)
}

func (node *BlockNode) GenerateNextBlock() *BlockNode {
	var nextBlockHeader = &BlockHeader{
		PrevBlock: node.hash,
		Version:   0,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
		Bits:      node.Bits,
	}

	nextBlockNode := newBlockNode(nextBlockHeader, node.Height+1)
	nextBlockNode.parent = node

	return nextBlockNode
}

func GetGenesisBlock() *BlockNode {
	zeroHash := chainhash.Hash{}
	var genesisHeader = &BlockHeader{
		PrevBlock: zeroHash,
		Version:   0,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
		Bits:      DefaultBits,
	}

	genesisNode := newBlockNode(genesisHeader, 0)

	return genesisNode
}
