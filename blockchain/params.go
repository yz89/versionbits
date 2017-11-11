package blockchain

const (
	// DefaultBits 默认bits，越小难度越大
	DefaultBits = uint32(0x20000009)
)

const (
	Miner1 uint32 = iota
	Miner2
	Miner3
	Miner4
)

// ConsensusDeployment defines details related to a specific consensus rule
// change that is voted in.  This is part of BIP0009.
type ConsensusDeployment struct {
	BitNumber  uint8
	StartTime  uint64
	ExpireTime uint64
}

// Constants that define the deployment offset in the deployments field of the
// parameters for each deployment.  This is useful to be able to get the details
// of a specific deployment by name.
const (
	DeploymentTestDummy = iota
	DeploymentCSV
	DeploymentSegwit
	DefinedDeployments
)

var (
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold = 1916 // 95% of MinerConfirmationWindow
	MinerConfirmationWindow       = 2016 //
	Deployments                   = []ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1462060800, // May 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1479168000, // November 15, 2016 UTC
			ExpireTime: 1510704000, // November 15, 2017 UTC.
		},
	}
	Deployments1 = []ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1462060800, // May 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
	}

	MinerDeploymentMap = map[uint32][]ConsensusDeployment{
		Miner1: Deployments,
		Miner2: Deployments,
		Miner3: Deployments1,
		Miner4: Deployments1,
	}
)
