package distributedchat

import (
	"net"

	"github.com/bits-and-blooms/bitset"
)

type raftNode struct {
	id               string
	nodeEpoch        uint64
	slots            *bitset.BitSet
	masterOfThisNode *raftNode

	conn net.Conn
}

type nodeInitParams struct {
	id      string
	address string
}

func (r *raftNode) init(params *nodeInitParams) error {
	r.id = params.id
	r.nodeEpoch = 0

	if len(params.address) > 0 {
		var err error
		r.conn, err = net.Dial("tcp", params.address)
		if err != nil {
			return err
		}
	}

	return nil
}
