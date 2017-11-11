package blockchain

const (
	// 001.......
	vbTopBits = 0x20000000
	// 111.......
	vbTopMask = 0xe0000000

	vbNumBits = 29
)

func GetNextVersion(minerID uint32) uint32 {
	expectVersion := uint32(vbTopBits)
	deployments := MinerDeploymentMap[minerID]
	for i := 0; i < len(deployments); i++ {
		expectVersion |= uint32(1) << deployments[i].BitNumber
	}
	return uint32(expectVersion)
}
