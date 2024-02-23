package distributedchat

import (
	"encoding/binary"
	"errors"
	"net"
	"pnas/prpc"

	"github.com/bits-and-blooms/bitset"
	"google.golang.org/protobuf/proto"
)

type raftNode struct {
	id               string
	actionEpoch      uint64
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
	r.actionEpoch = 0

	if len(params.address) > 0 {
		var err error
		r.conn, err = net.Dial("tcp", params.address)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *raftNode) send(msg *prpc.RaftMsg) error {
	if r.conn == nil {
		return errors.New("nil conn")
	}
	msgBuf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	var buf []byte
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(msgBuf)))
	buf = append(buf, msgBuf...)
	_, err = r.conn.Write(buf)
	return err
}
