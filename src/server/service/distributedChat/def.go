package distributedchat

type NodeRole int

const (
	SlotsNum int = 1 << 14

	MasterRole NodeRole = 1
	SlaveRole  NodeRole = 2
)
