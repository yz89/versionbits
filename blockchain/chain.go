package blockchain

import "versionbits/chainhash"

type BlockNode struct {
	parent    *BlockNode
	hash      chainhash.Hash
	Height    uint32
	Bits      uint32
	Version   uint32
	Nonce     uint32
	Timestamp int64
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
