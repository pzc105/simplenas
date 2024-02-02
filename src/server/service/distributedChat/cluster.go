package distributedchat

import (
	"context"
	"log"
	"net"
	"pnas/prpc"
	"sync"
	"sync/atomic"

	"github.com/bits-and-blooms/bitset"
	"google.golang.org/grpc"
)

type Cluster struct {
	prpc.NodeServiceServer

	mtx        sync.Mutex
	epoch      uint64
	nodes      map[string]*raftNode
	slots2Node []*raftNode

	connMtx   sync.Mutex
	conns     map[string]net.Conn
	conn2Node map[string]*raftNode

	shutdown    atomic.Bool
	shutdownCtx context.Context
	closeFunc   context.CancelFunc
	wg          sync.WaitGroup

	grpcSer *grpc.Server
}

type JoinParams struct {
	ListenAddress string
	TargetAddress string
	TargetId      string
	Role          NodeRole
	Slots         *bitset.BitSet
}

func (c *Cluster) Start(params *JoinParams) {
	c.shutdownCtx, c.closeFunc = context.WithCancel(context.Background())
	c.shutdown.Store(false)

	if len(params.TargetAddress) > 0 {
		c.wg.Add(1)
		go func() {
			defer c.wg.Add(-1)
			rn := &raftNode{}
			rn.init(&nodeInitParams{
				id:      params.TargetId,
				address: params.TargetAddress,
			})
			c.mtx.Lock()
			c.nodes[params.TargetId] = rn
			c.mtx.Unlock()

		}()
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
		c.connMtx.Lock()
		c.conns[conn.RemoteAddr().String()] = conn
		c.connMtx.Unlock()
	}
}
