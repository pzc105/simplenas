package distributedchat

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"pnas/prpc"
	"pnas/utils"
	"sync"
	"sync/atomic"

	"github.com/bits-and-blooms/bitset"
	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	mtx               sync.Mutex
	nextVoteEpoch     atomic.Uint64
	nextActionEpoch   atomic.Uint64
	lastCommitedEpoch atomic.Uint64
	actors            map[string]*raftNode
	nodes             map[string]*raftNode
	myself            *raftNode
	master            *raftNode
	slots2Node        []*raftNode

	isMaster atomic.Bool

	shutdown    atomic.Bool
	shutdownCtx context.Context
	closeFunc   context.CancelFunc
	wg          sync.WaitGroup
}

type StartParams struct {
	ListenAddress string
	TargetAddress string
	MyId          string
	TargetId      string
	Role          NodeRole
	Slots         *bitset.BitSet
}

func (c *Cluster) Start(params *StartParams) {
	c.shutdownCtx, c.closeFunc = context.WithCancel(context.Background())
	c.shutdown.Store(false)

	c.nextActionEpoch.Store(1)
	c.nextVoteEpoch.Store(1)
	c.lastCommitedEpoch.Store(0)
	c.slots2Node = make([]*raftNode, SlotsNum)

	c.myself = &raftNode{}
	c.myself.init(&nodeInitParams{
		id: params.MyId,
	})
	c.actors = make(map[string]*raftNode)
	if params.Role == ActorRole {
		c.actors[params.MyId] = c.myself
	}

	c.nodes = make(map[string]*raftNode)
	c.nodes[params.MyId] = c.myself

	if len(params.TargetAddress) > 0 {
		c.isMaster.Store(false)
		c.wg.Add(1)
		go func() {
			defer c.wg.Add(-1)
			rn := &raftNode{}
			err := rn.init(&nodeInitParams{
				id:      params.TargetId,
				address: params.TargetAddress,
			})
			if err != nil {
				panic(err)
			}
			c.mtx.Lock()
			c.nodes[params.TargetId] = rn
			c.mtx.Unlock()
			var joinMsg prpc.RaftMsg
			joinMsg.Type = prpc.RaftMsg_Action
			joinMsg.Action = &prpc.RaftTransaction{
				Type: prpc.RaftTransaction_NewNode,
				NewNode: &prpc.NewNodeAction{
					MyId:      params.MyId,
					MyAddress: params.ListenAddress,
					Role:      int32(params.Role),
				},
			}
			rn.send(&joinMsg)
		}()
	} else {
		c.isMaster.Store(true)
		c.master = c.myself
	}

	lis, err := net.Listen("tcp", params.ListenAddress)
	if err != nil {
		log.Panic(err)
	}

	for !c.shutdown.Load() {
		conn, err := lis.Accept()
		if err != nil {
			log.Panic(err)
		}
		c.wg.Add(1)
		go func() {
			defer c.wg.Add(-1)
			c.handleConnRead(conn)
		}()
	}
}

func (c *Cluster) handleConnRead(conn net.Conn) {
	pendingBuf := make([]byte, 1024)
	var buf []byte
	for {
		l, err := conn.Read(pendingBuf)
		if err != nil {
			return
		}
		buf = append(buf, pendingBuf[:l]...)
		if len(buf) >= 4 {
			dataLen := (int)(binary.BigEndian.Uint32(buf))
			if dataLen+4 <= len(buf) {
				oneMsgBuf := buf[4 : 4+dataLen]
				buf = buf[4+dataLen:]
				var msg prpc.RaftMsg
				err = proto.Unmarshal(oneMsgBuf, &msg)
				if err == nil {
					c.handleMsg(&msg)
				} else {
					conn.Close()
					return
				}
			}
		}
	}
}

func (c *Cluster) handleMsg(msg *prpc.RaftMsg) {
	c.mtx.Lock()
	master := c.master
	c.mtx.Unlock()

	if msg.Type == prpc.RaftMsg_Action {
		if c.isMaster.Load() {
			c.startAction(msg)
		} else if msg.Action.MyId != master.id {
			master.send(msg)
		} else {

		}
		return
	}
}

func (c *Cluster) startAction(msg *prpc.RaftMsg) {
	newepoch := utils.FetchAndAdd(&c.nextActionEpoch, uint64(1))
	msg.Action.Epoch = newepoch
	msg.Action.MyId = c.myself.id

}

func ()