package blockchain

import (
	"fmt"
)

func thresholdState(prevNode *BlockNode) {
	// confirmationWindow := MinerConfirmationWindow
	var neededStates []*BlockNode
	for prevNode != nil {
		neededStates = append(neededStates, prevNode)

		prevNode = prevNode.Parent
	}

	for neededNum := len(neededStates) - 1; neededNum >= 0; neededNum-- {
		prevNode := neededStates[neededNum]

		fmt.Print(" ", &prevNode)
	}
	fmt.Println()
}
