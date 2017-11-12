package blockchain

const (
	// 001.......
	vbTopBits = 0x20000000
	// 111.......
	vbTopMask = 0xe0000000

	vbNumBits = 29
)

type deploymentChecker struct {
	deployment *ConsensusDeployment
}

var _ thresholdConditionChecker = deploymentChecker{}

func (c deploymentChecker) BeginTime() uint64 {
	return c.deployment.StartTime
}

func (c deploymentChecker) EndTime() uint64 {
	return c.deployment.ExpireTime
}

func (c deploymentChecker) RuleChangeActivationThreshold() uint32 {
	return uint32(RuleChangeActivationThreshold)
}

func (c deploymentChecker) MinerConfirmationWindow() uint32 {
	return uint32(MinerConfirmationWindow)
}

func (c deploymentChecker) Condition(node *BlockNode) (bool, error) {
	conditionMask := uint32(1) << c.deployment.BitNumber
	version := uint32(node.Version)
	return (version&vbTopMask == vbTopBits) && (version&conditionMask != 0), nil
}

// GetNextVersion 根据矿工deployment配置，获取下一个块的version值
// 同时也是bip9实现的函数入口，里面会推导状态并缓存
func GetNextVersion(minerID uint32, prevNode *BlockNode) uint32 {
	expectVersion := uint32(vbTopBits)
	deployments := MinerDeploymentMap[minerID]
	for i := 0; i < len(deployments); i++ {
		deployment := &deployments[i]
		checker := deploymentChecker{deployment: deployment}
		thresholdState(prevNode, checker)
		expectVersion |= uint32(1) << deployments[i].BitNumber
	}
	return uint32(expectVersion)
}
