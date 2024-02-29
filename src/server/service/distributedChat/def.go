package distributedchat

type NodeRole int

type NodeState int

type ClusterState int

const (
	MasterRole NodeRole = 1
	ActorRole  NodeRole = 2
	SlaveRole  NodeRole = 3
)

const (
	UnknownNodeState   NodeState = 0
	MaybeDownNodeState NodeState = 1
	DownNodeState      NodeState = 2
	ConnectedNodeState NodeState = 3
)

const (
	NormalCluster   ClusterState = 1
	ElectionCluster ClusterState = 2
)

const (
	HashSlotSize = 1 << 14

	SyncSlotDataStep = 1
	SwitchSlotStep   = 2
)
