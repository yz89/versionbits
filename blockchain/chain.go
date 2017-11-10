package blockchain

import (
	"sync"
)

type Chain struct {
	stateLock sync.RWMutex
	tipBlock  *BlockNode
}

func (c *Chain) Tip() *BlockNode {
	c.stateLock.Lock()
	tipBlock := c.tipBlock
	c.stateLock.Unlock()

	return tipBlock
}

func (c *Chain) AddBest(bestBlock *BlockNode) bool {
	ok := false
	c.stateLock.Lock()
	if bestBlock.Parent == nil || (bestBlock.Parent.Hash == c.tipBlock.Hash && bestBlock.Height == c.tipBlock.Height+1) {
		c.tipBlock = bestBlock
		ok = true
	}
	c.stateLock.Unlock()
	return ok
}
