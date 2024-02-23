package distributedchat

type NodeRole int

type NodeState int

const (
	SlotsNum int = 1 << 14

	MasterRole NodeRole = 1
	ActorRole  NodeRole = 2
	SlaveRole  NodeRole = 3
)
