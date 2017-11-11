package blockchain

import (
	"math/rand"
	"time"
	"versionbits/chainhash"
)

type BlockNode struct {
	Parent     *BlockNode
	Hash       chainhash.Hash
	MerkleRoot chainhash.Hash
	Height     uint32
	Version    uint32
	Timestamp  int64
	Nonce      uint32
	Bits       uint32
}

func initBlockNode(node *BlockNode, blockHeader *BlockHeader, height uint32) {
	*node = BlockNode{
		Hash:       blockHeader.HashBlock(),
		MerkleRoot: blockHeader.MerkleRoot,
		Height:     height,
		Version:    blockHeader.Version,
		Timestamp:  blockHeader.Timestamp,
		Nonce:      blockHeader.Nonce,
		Bits:       blockHeader.Bits,
	}
}

func newBlockNode(blockHeader *BlockHeader, height uint32) *BlockNode {
	var node BlockNode
	initBlockNode(&node, blockHeader, height)
	return &node
}

func (node *BlockNode) Header() *BlockHeader {
	prevHash := chainhash.Hash{}
	if node.Parent != nil {
		prevHash = node.Parent.Hash
	}
	return &BlockHeader{
		Version:    node.Version,
		PrevBlock:  prevHash,
		MerkleRoot: node.MerkleRoot,
		Timestamp:  node.Timestamp,
		Bits:       node.Bits,
		Nonce:      node.Nonce,
	}
}

func (node *BlockNode) Ancestor(height uint32) *BlockNode {
	if height < 0 || height > node.Height {
		return nil
	}

	n := node
	for ; n != nil && n.Height != height; n = n.Parent {
		// Intentionally left blank
	}
	return n
}

func (node *BlockNode) RelativeAncestor(distance uint32) *BlockNode {
	return node.Ancestor(node.Height - distance)
}

func (node *BlockNode) GenerateNextBlock() *BlockNode {
	// simulate different transactions collection
	transactions := []byte{byte(rand.Intn(127)), byte(rand.Intn(127)), byte(rand.Intn(127)), byte(rand.Intn(127))}

	// calculate next difficulty
	bits := CalcNextRequiredDifficulty(node)

	var nextBlockHeader = &BlockHeader{
		PrevBlock:  node.Hash,
		MerkleRoot: chainhash.HashH(transactions),
		Version:    0,
		Timestamp:  time.Now().Unix(),
		Nonce:      0,
		Bits:       bits,
	}

	nextBlockNode := newBlockNode(nextBlockHeader, node.Height+1)
	nextBlockNode.Parent = node

	return nextBlockNode
}

func GetGenesisBlock() *BlockNode {
	zeroHash := chainhash.Hash{}
	var genesisHeader = &BlockHeader{
		PrevBlock:  zeroHash,
		MerkleRoot: zeroHash,
		Version:    0,
		Timestamp:  time.Now().Unix(),
		Nonce:      0,
		Bits:       DefaultBits,
	}

	genesisNode := newBlockNode(genesisHeader, 0)
	genesisNode.Parent = nil

	return genesisNode
}
