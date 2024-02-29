package distributedchat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"pnas/prpc"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

type raftNode struct {
	mtx                 sync.Mutex
	id                  string
	states              map[string]NodeState
	currentEpoch        atomic.Int64
	lastCommitedEpoch   atomic.Int64
	role                NodeRole
	address             string
	lastVoteEpoch       int64
	voteId              string
	responsibleSlots    []byte
	lastSyncActionsTime time.Time

	conn net.Conn
}

type nodeInitParams struct {
	id      string
	address string
	role    NodeRole
}

func (n *raftNode) init(params *nodeInitParams) error {
	n.id = params.id
	n.currentEpoch.Store(0)
	n.lastCommitedEpoch.Store(0)
	n.role = params.role
	n.states = make(map[string]NodeState)
	n.address = params.address
	n.responsibleSlots = make([]byte, HashSlotSize/8)
	n.lastSyncActionsTime = time.Now()

	if len(params.address) > 0 {
		var err error
		fmt.Printf("connecting:%s, id:%s, role:%d\n", params.address, params.id, params.role)
		n.conn, err = net.Dial("tcp", params.address)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *raftNode) send(msg *prpc.RaftMsg) error {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	if n.conn == nil {
		n.reconnect()
		if n.conn == nil {
			return errors.New("nil conn")
		}
	}
	msgBuf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	var buf []byte
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(msgBuf)))
	buf = append(buf, msgBuf...)
	retryCount := 0
retry:
	_, err = n.conn.Write(buf)
	if err != nil && retryCount < 2 {
		retryCount++
		err = n.reconnect()
		if err == nil {
			goto retry
		}
	}
	return err
}

func (n *raftNode) reconnect() error {
	if n.conn != nil {
		n.conn.Close()
	}
	if len(n.address) > 0 {
		var err error
		n.conn, err = net.Dial("tcp", n.address)
		if err == nil {
			fmt.Printf("reconnect to:%s\n", n.id)
		}
		return err
	} else {
		return errors.New("")
	}
}

func (n *raftNode) changeCurrentEpoch(epoch int64) {
	if epoch == n.currentEpoch.Load() {
		return
	}
	fmt.Printf("[epoch] change id:%s current epoch:%d->%d\n", n.id, n.currentEpoch.Load(), epoch)
	n.currentEpoch.Store(epoch)
}

func (n *raftNode) changeCommitedEpoch(epoch int64) {
	if epoch == n.lastCommitedEpoch.Load() {
		return
	}
	fmt.Printf("[epoch] change id:%s commited epoch:%d->%d\n", n.id, n.lastCommitedEpoch.Load(), epoch)
	n.lastCommitedEpoch.Store(epoch)
}

func (n *raftNode) updateOtherNodeState(id string, s NodeState) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	n.states[id] = s
}

func (n *raftNode) getOtherNodeState(id string) NodeState {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	if s, ok := n.states[id]; ok {
		return s
	}
	return UnknownNodeState
}

func (n *raftNode) isSet(i int) bool {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	return isSet(n.responsibleSlots, i)
}

func (n *raftNode) setSlot(i int) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	setSlot(n.responsibleSlots, i)
}

func (n *raftNode) unsetSlot(i int) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	unsetSlot(n.responsibleSlots, i)
}
