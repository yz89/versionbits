package blockchain

import (
	"math/rand"
	"sort"
	"time"
	"versionbits/chainhash"
)

type BlockNode struct {
	Parent     *BlockNode
	Hash       chainhash.Hash
	MerkleRoot chainhash.Hash
	Height     int32
	Version    uint32
	Timestamp  int64
	Nonce      uint32
	Bits       uint32
}

func initBlockNode(node *BlockNode, blockHeader *BlockHeader, height int32) {
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

func newBlockNode(blockHeader *BlockHeader, height int32) *BlockNode {
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

func (node *BlockNode) Ancestor(height int32) *BlockNode {
	if height < 0 || height > node.Height {
		return nil
	}

	n := node
	for ; n != nil && n.Height != height; n = n.Parent {
		// Intentionally left blank
	}
	return n
}

func (node *BlockNode) RelativeAncestor(distance int32) *BlockNode {
	return node.Ancestor(node.Height - distance)
}

func (node *BlockNode) CalcPastMedianTime() time.Time {
	// Create a slice of the previous few block timestamps used to calculate
	// the median per the number defined by the constant medianTimeBlocks.
	timestamps := make([]int64, medianTimeBlocks)
	numNodes := 0
	iterNode := node
	for i := 0; i < medianTimeBlocks && iterNode != nil; i++ {
		timestamps[i] = iterNode.Timestamp
		numNodes++

		iterNode = iterNode.Parent
	}

	// Prune the slice to the actual number of available timestamps which
	// will be fewer than desired near the beginning of the block chain
	// and sort them.
	timestamps = timestamps[:numNodes]
	sort.Sort(timeSorter(timestamps))

	medianTimestamp := timestamps[numNodes/2]
	return time.Unix(medianTimestamp, 0)
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
