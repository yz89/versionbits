package blockchain

const (
	// 001.......
	vbTopBits = 0x20000000
	// 111.......
	vbTopMask = 0xe0000000

	vbNumBits = 29
)

// GetNextVersion 根据矿工deployment配置，获取下一个块的version值
// 同时也是bip9实现的函数入口，里面会推导状态并缓存
func GetNextVersion(minerID uint32, prevNode *BlockNode) uint32 {
	expectVersion := uint32(vbTopBits)
	deployments := MinerDeploymentMap[minerID]
	for i := 0; i < len(deployments); i++ {
		expectVersion |= uint32(1) << deployments[i].BitNumber
	}
	thresholdState(prevNode)
	return uint32(expectVersion)
}
