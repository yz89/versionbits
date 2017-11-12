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

// thresholdConditionChecker provides a generic interface that is invoked to
// determine when a consensus rule change threshold should be changed.
type thresholdConditionChecker interface {
	// BeginTime returns the unix timestamp for the median block time after
	// which voting on a rule change starts (at the next window).
	BeginTime() uint64

	// EndTime returns the unix timestamp for the median block time after
	// which an attempted rule change fails if it has not already been
	// locked in or activated.
	EndTime() uint64

	// RuleChangeActivationThreshold is the number of blocks for which the
	// condition must be true in order to lock in a rule change.
	RuleChangeActivationThreshold() uint32

	// MinerConfirmationWindow is the number of blocks in each threshold
	// state retarget window.
	MinerConfirmationWindow() uint32

	// Condition returns whether or not the rule change activation condition
	// has been met.  This typically involves checking whether or not the
	// bit assocaited with the condition is set, but can be more complex as
	// needed.
	Condition(*BlockNode) (bool, error)
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

func thresholdState(prevNode *BlockNode, checker thresholdConditionChecker) (ThresholdState, error) {

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

		switch state {
		case ThresholdDefined:
			// The deployment of the rule change fails if it expires
			// before it is accepted and locked in.
			medianTime := prevNode.CalcPastMedianTime()
			medianTimeUnix := uint64(medianTime.Unix())
			if medianTimeUnix >= checker.EndTime() {
				state = ThresholdFailed
				break
			}

			// The state for the rule moves to the started state
			// once its start time has been reached (and it hasn't
			// already expired per the above).
			if medianTimeUnix >= checker.BeginTime() {
				state = ThresholdStarted
			}

		case ThresholdStarted:
			// The deployment of the rule change fails if it expires
			// before it is accepted and locked in.
			medianTime := prevNode.CalcPastMedianTime()
			if uint64(medianTime.Unix()) >= checker.EndTime() {
				state = ThresholdFailed
				break
			}

			// At this point, the rule change is still being voted
			// on by the miners, so iterate backwards through the
			// confirmation window to count all of the votes in it.
			var count uint32
			countNode := prevNode
			for i := int32(0); i < confirmationWindow; i++ {
				condition, err := checker.Condition(countNode)
				if err != nil {
					return ThresholdFailed, err
				}
				if condition {
					count++
				}

				// Get the previous block node.
				countNode = countNode.Parent
			}

			// The state is locked in if the number of blocks in the
			// period that voted for the rule change meets the
			// activation threshold.
			if count >= checker.RuleChangeActivationThreshold() {
				state = ThresholdLockedIn
			}

		case ThresholdLockedIn:
			// The new rule becomes active when its previous state
			// was locked in.
			state = ThresholdActive

		// Nothing to do if the previous state is active or failed since
		// they are both terminal states.
		case ThresholdActive:
		case ThresholdFailed:
		}
	}
	return state, nil
}
