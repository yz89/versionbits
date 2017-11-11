package blockchain

import (
	"fmt"
	"versionbits/chainhash"
)

type ThresholdState byte

// These constants are used to identify specific threshold states.
const (
	ThresholdDefined ThresholdState = iota
	ThresholdStarted
	ThresholdLockedIn
	ThresholdActive
	ThresholdFailed
	numThresholdsStates
)

// thresholdStateStrings is a map of ThresholdState values back to their
// constant names for pretty printing.
var thresholdStateStrings = map[ThresholdState]string{
	ThresholdDefined:  "ThresholdDefined",
	ThresholdStarted:  "ThresholdStarted",
	ThresholdLockedIn: "ThresholdLockedIn",
	ThresholdActive:   "ThresholdActive",
	ThresholdFailed:   "ThresholdFailed",
}

// String returns the ThresholdState as a human-readable name.
func (t ThresholdState) String() string {
	if s := thresholdStateStrings[t]; s != "" {
		return s
	}
	return fmt.Sprintf("Unknown ThresholdState (%d)", int(t))
}

// thresholdStateCache provides a type to cache the threshold states of each
// threshold window for a set of IDs.
type thresholdStateCache struct {
	entries map[chainhash.Hash]ThresholdState
}

// Lookup returns the threshold state associated with the given hash along with
// a boolean that indicates whether or not it is valid.
func (c *thresholdStateCache) Lookup(hash *chainhash.Hash) (ThresholdState, bool) {
	state, ok := c.entries[*hash]
	return state, ok
}

// Update updates the cache to contain the provided hash to threshold state
// mapping.
func (c *thresholdStateCache) Update(hash *chainhash.Hash, state ThresholdState) {
	c.entries[*hash] = state
}

// newThresholdCaches returns a new array of caches to be used when calculating
// threshold states.
func newThresholdCaches(numCaches uint32) []thresholdStateCache {
	caches := make([]thresholdStateCache, numCaches)
	for i := 0; i < len(caches); i++ {
		caches[i] = thresholdStateCache{
			entries: make(map[chainhash.Hash]ThresholdState),
		}
	}
	return caches
}

func thresholdState(prevNode *BlockNode) (ThresholdState, error) {

	confirmationWindow := int32(MinerConfirmationWindow)
	if prevNode == nil || (prevNode.Height+1) < confirmationWindow {
		return ThresholdDefined, nil
	}

	prevNode = prevNode.Ancestor(prevNode.Height -
		(prevNode.Height+1)%confirmationWindow)

	var neededStates []*BlockNode
	for prevNode != nil {
		neededStates = append(neededStates, prevNode)

		prevNode = prevNode.RelativeAncestor(confirmationWindow)
	}

	state := ThresholdDefined
	for neededNum := len(neededStates) - 1; neededNum >= 0; neededNum-- {
		prevNode := neededStates[neededNum]

		fmt.Printf("%d ", prevNode.Height)
	}
	fmt.Println()
	return state, nil
}
