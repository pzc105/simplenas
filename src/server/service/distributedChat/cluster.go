package distributedchat

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"pnas/prpc"
	"pnas/utils"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	mtx              sync.Mutex
	actors           map[string]*raftNode
	nodes            map[string]*raftNode
	myself           *raftNode
	master           *raftNode
	slots2Node       []*raftNode
	actions          map[int64]*prpc.RaftTransaction
	voteEpoch        atomic.Int64
	nextVoteEpoch    int64
	isMaster         atomic.Bool
	lastAckTime      time.Time
	maxMasterTimeout time.Duration
	state            ClusterState
	rand             *rand.Rand
	externalAddress  string
	pendingSlots     []byte

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
	StartSlot     int
	EndSlot       int
}

func (c *Cluster) Wait() {
	c.wg.Wait()
}

func (c *Cluster) Start(params *StartParams) {
	c.shutdownCtx, c.closeFunc = context.WithCancel(context.Background())
	c.shutdown.Store(false)

	c.slots2Node = make([]*raftNode, HashSlotSize)
	c.actions = make(map[int64]*prpc.RaftTransaction)

	c.myself = &raftNode{}
	c.myself.init(&nodeInitParams{
		id:   params.MyId,
		role: params.Role,
	})
	c.actors = make(map[string]*raftNode)
	if params.Role == ActorRole {
		c.actors[params.MyId] = c.myself
	}

	c.nodes = make(map[string]*raftNode)
	c.nodes[params.MyId] = c.myself
	c.voteEpoch.Store(0)
	c.nextVoteEpoch = 1
	c.maxMasterTimeout = time.Second * 3
	c.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	c.externalAddress = params.ListenAddress

	c.pendingSlots = make([]byte, HashSlotSize/8)
	fmt.Printf("slot:%d->%d\n", params.StartSlot, params.EndSlot)
	for i := params.StartSlot; i <= params.EndSlot; i++ {
		if i < 0 || i >= HashSlotSize {
			continue
		}
		setSlot(c.pendingSlots, i)
	}

	if len(params.TargetAddress) > 0 {
		c.join(params.TargetId, params.TargetAddress)
	} else {
		c.isMaster.Store(true)
		c.voteEpoch.Store(1)
		c.nextVoteEpoch = 2
		c.master = c.myself
		c.myself.role = MasterRole
		action := prpc.RaftTransaction{
			MyId:  c.myself.id,
			Type:  prpc.RaftTransaction_NewNode,
			Epoch: 1,
			NewNode: &prpc.NewNodeAction{
				MyId:      c.myself.id,
				MyAddress: params.ListenAddress,
				Role:      int32(ActorRole),
			},
		}
		action2 := prpc.RaftTransaction{
			MyId:  c.myself.id,
			Type:  prpc.RaftTransaction_HashSlotAction,
			Epoch: 2,
			HashSlot: &prpc.HashSlotAction{
				MyId: c.myself.id,
				Step: SwitchSlotStep,
			},
		}
		action2.HashSlot.Slots = make([]byte, HashSlotSize/8)
		copy(action2.HashSlot.Slots, c.pendingSlots)
		c.myself.responsibleSlots = c.pendingSlots
		for i := 0; i < HashSlotSize; i++ {
			if c.myself.isSet(i) {
				c.slots2Node[i] = c.myself
			}
		}
		c.pendingSlots = []byte{}
		c.actions[1] = &action
		c.actions[2] = &action2
		c.myself.currentEpoch.Store(2)
		c.myself.lastCommitedEpoch.Store(2)
	}

	fmt.Printf("myid:%s, role:%d, listen:%s\n", params.MyId, params.Role, params.ListenAddress)
	lis, err := net.Listen("tcp", params.ListenAddress)
	if err != nil {
		log.Panic(err)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Add(-1)
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
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Add(-1)
		tick := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-tick.C:
				if !c.isMaster.Load() {
					c.mtx.Lock()
					master := c.master
					lastAckTime := c.lastAckTime
					if master == nil {
						c.mtx.Unlock()
						continue
					}
					master.send(c.newPingMsg())
					if time.Since(lastAckTime) >= c.maxMasterTimeout {
						c.markNodeMaybeDownLocked(master.id)
					}

					if c.master != nil && len(c.pendingSlots) == HashSlotSize/8 {
						action := prpc.RaftTransaction{
							MyId: c.myself.id,
							Type: prpc.RaftTransaction_HashSlotAction,
							HashSlot: &prpc.HashSlotAction{
								MyId: c.myself.id,
								Step: SwitchSlotStep,
							},
						}
						action.HashSlot.Slots = make([]byte, HashSlotSize/8)
						copy(action.HashSlot.Slots, c.pendingSlots)
						msg := prpc.RaftMsg{
							Type:   prpc.RaftMsg_Action,
							Action: &action,
						}
						c.pendingSlots = []byte{}
						c.startActionLocked(&msg)
					}
					c.mtx.Unlock()
				}
			case <-c.shutdownCtx.Done():
				return
			}
		}
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Add(-1)
		tick := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-tick.C:
				c.mtx.Lock()
				nodes := make([]*raftNode, 0, len(c.nodes))
				for _, n := range c.nodes {
					nodes = append(nodes, n)
				}
				randNode := c.getRandOtherNodeLocked()
				slot := c.rand.Int() % HashSlotSize
				c.mtx.Unlock()

				pingMsg := c.newPingMsg()
				for _, n := range nodes {
					if n.id == c.myself.id {
						continue
					}
					n.send(pingMsg)
				}
				msg := fmt.Sprintf("from:%s, now:%s", c.myself.id, time.Now().String())
				if randNode != nil {
					randNode.send(&prpc.RaftMsg{
						Type: prpc.RaftMsg_SendMsg2Slot,
						SlotMsg: &prpc.SlotMsg{
							Slot: int32(slot),
							Msg:  msg,
						},
					})
				}

			case <-c.shutdownCtx.Done():
				return
			}
		}
	}()
}

func (c *Cluster) join(targetId string, targetAddress string) {
	c.isMaster.Store(false)
	c.wg.Add(1)
	go func() {
		defer c.wg.Add(-1)
		rn := &raftNode{}
		err := rn.init(&nodeInitParams{
			id:      targetId,
			address: targetAddress,
		})
		if err != nil {
			panic(err)
		}
		c.mtx.Lock()
		c.nodes[targetId] = rn
		c.lastAckTime = time.Now()
		c.mtx.Unlock()
		var joinMsg prpc.RaftMsg
		joinMsg.Type = prpc.RaftMsg_Action
		joinMsg.Action = &prpc.RaftTransaction{
			Type: prpc.RaftTransaction_NewNode,
			NewNode: &prpc.NewNodeAction{
				MyId:      c.myself.id,
				MyAddress: c.externalAddress,
				Role:      int32(c.myself.role),
			},
		}
		rn.send(&joinMsg)
	}()
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
		for len(buf) >= 4 {
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
			} else {
				break
			}
		}
	}
}

func (c *Cluster) handleMsg(msg *prpc.RaftMsg) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if msg.Type == prpc.RaftMsg_Action {
		if c.master == nil {
			return
		}
		c.startActionLocked(msg)
		return
	} else if msg.Type == prpc.RaftMsg_SyncAction {
		c.handleSyncActionsLocked(msg)
	} else if msg.Type == prpc.RaftMsg_SyncActionRet {
		c.handleSyncActionsRetLocked(msg)
	} else if msg.Type == prpc.RaftMsg_Ping {
		c.handlePingLocked(msg)
	} else if msg.Type == prpc.RaftMsg_Pong {
		c.handlePongMsgLock(msg)
	} else if msg.Type == prpc.RaftMsg_Election {
		c.handleElectionMsgLocked(msg)
	} else if msg.Type == prpc.RaftMsg_ElectionRet {
		c.handleElectionRetMsgLocked(msg)
	} else if msg.Type == prpc.RaftMsg_ReqActions {
		c.handleReqActionsLocked(msg)
	} else if msg.Type == prpc.RaftMsg_SendMsg2Slot {
		c.handleSoltMsgLocked(msg)
	}
}

func (c *Cluster) handleSyncActionsLocked(msg *prpc.RaftMsg) {
	if c.isMaster.Load() {
		return
	}
	syncMsg := msg.SyncActions
	if len(syncMsg.Actions) == 0 {
		return
	}
	minEpoch := int64(-1)
	maxEpoch := int64(-1)
	for _, a := range syncMsg.Actions {
		if a.Epoch > c.myself.lastCommitedEpoch.Load() {
			c.actions[a.Epoch] = a
			if a.Epoch > maxEpoch || maxEpoch == -1 {
				maxEpoch = a.Epoch
			}
			if a.Epoch < minEpoch || minEpoch == -1 {
				minEpoch = a.Epoch
			}
		}
	}

	if maxEpoch == -1 {
		return
	}

	fmt.Printf("[sync action] minEpoch:%d maxEpoch:%d\n", minEpoch, maxEpoch)

	if minEpoch > c.myself.lastCommitedEpoch.Load()+1 {
		fmt.Printf("[action] need request old action. recv minEpoch:%d\n", minEpoch)
		if c.master != nil {
			c.master.send(c.newPingMsg())
		}
		return
	}
	c.myself.changeCurrentEpoch(maxEpoch)
	if syncMsg.CommitedEpoch > c.myself.lastCommitedEpoch.Load() {
		c.handleCommitedActionLocked(syncMsg.CommitedEpoch)
	}
	if c.master != nil {
		retMsg := &prpc.RaftMsg{
			Type: prpc.RaftMsg_SyncActionRet,
			SyncActionsRet: &prpc.RaftSyncActionsRet{
				MyId:         c.myself.id,
				CurrentEpoch: c.myself.currentEpoch.Load(),
			},
		}
		c.master.send(retMsg)
	}
}

func (c *Cluster) handleSyncActionsRetLocked(msg *prpc.RaftMsg) {
	retMsg := msg.SyncActionsRet
	if c.isMaster.Load() && retMsg.MyId != c.myself.id {
		c.addNodeEpochLocked(retMsg.MyId, retMsg.CurrentEpoch)
	}
}

func (c *Cluster) handleReqActionsLocked(msg *prpc.RaftMsg) {
	req := msg.ReqActions
	node, ok := c.nodes[req.MyId]
	if !ok {
		return
	}
	node.lastCommitedEpoch.Store(req.CommitedEpoch)
	c.syncActionsLocked(node)
}

func (c *Cluster) handleSoltMsgLocked(msg *prpc.RaftMsg) {
	slotMsg := msg.SlotMsg
	if c.myself.isSet(int(slotMsg.Slot)) {
		fmt.Printf("recv msg, slot:%d msg:%s\n", slotMsg.Slot, slotMsg.Msg)
	} else if c.isMaster.Load() {
		node := c.slots2Node[slotMsg.Slot]
		if node != nil {
			node.send(msg)
		}
	} else {
		if c.master != nil {
			c.master.send(msg)
		}
	}
}

func (c *Cluster) newPingMsg() *prpc.RaftMsg {
	msg := prpc.RaftMsg{
		Type: prpc.RaftMsg_Ping,
		Ping: &prpc.RaftPing{
			MyId:          c.myself.id,
			Role:          int32(c.myself.role),
			CurrentEpoch:  c.myself.currentEpoch.Load(),
			CommitedEpoch: c.myself.lastCommitedEpoch.Load(),
			VoteEpoch:     c.voteEpoch.Load(),
		},
	}
	if c.isMaster.Load() {
		msg.Ping.MasterAddress = c.externalAddress
	}
	return &msg
}

func (c *Cluster) newPongMsg() *prpc.RaftMsg {
	msg := prpc.RaftMsg{
		Type: prpc.RaftMsg_Pong,
		Pong: &prpc.RaftPong{
			MyId: c.myself.id,
			Role: int32(c.myself.role),
		},
	}
	return &msg
}

func (c *Cluster) startActionLocked(msg *prpc.RaftMsg) {
	if c.isMaster.Load() {
		newepoch := utils.FetchAndAdd(&c.myself.currentEpoch, int64(1)) + 1
		msg.Action.Epoch = newepoch
		msg.Action.MyId = c.myself.id
		c.actions[newepoch] = msg.Action

		fmt.Printf("new action epoch:%d type:%d\n", msg.Action.Epoch, msg.Action.Type)

		c.addNodeEpochLocked(c.myself.id, newepoch)
		for _, n := range c.nodes {
			c.syncActionsLocked(n)
		}
	} else if c.master != nil {
		c.master.send(msg)
	}
}

func (c *Cluster) syncActionsLocked(node *raftNode) {
	if node == c.myself {
		return
	}
	epoch := node.lastCommitedEpoch.Load() + 1
	msg := prpc.RaftMsg{
		Type: prpc.RaftMsg_SyncAction,
		SyncActions: &prpc.RaftSyncActions{
			MyId:          c.myself.id,
			CurrentEpoch:  c.myself.currentEpoch.Load(),
			CommitedEpoch: c.myself.lastCommitedEpoch.Load(),
		},
	}
	for ; epoch <= c.myself.currentEpoch.Load(); epoch++ {
		if a, ok := c.actions[epoch]; ok {
			msg.SyncActions.Actions = append(msg.SyncActions.Actions, a)
		}
	}
	if len(msg.SyncActions.Actions) > 0 {
		err := node.send(&msg)
		if err == nil {
			node.lastSyncActionsTime = time.Now()
		}
		startEpoch := msg.SyncActions.Actions[0].Epoch
		endEpoch := msg.SyncActions.Actions[len(msg.SyncActions.Actions)-1].Epoch
		fmt.Printf("sync actions to id:%s, start epoch:%d, end epoch:%d, err:%v\n", node.id, startEpoch, endEpoch, err)
	}
}

func (c *Cluster) handlePingLocked(msg *prpc.RaftMsg) {
	pingMsg := msg.Ping
	if pingMsg.MyId == c.myself.id {
		return
	}

	c.handleCommitedActionLocked(pingMsg.CommitedEpoch)

	node, ok := c.nodes[pingMsg.MyId]
	if !ok {
		fmt.Printf("not found node. from:%s vote epoch:%d role:%d\n", pingMsg.MyId, pingMsg.VoteEpoch, pingMsg.Role)
		if pingMsg.VoteEpoch > c.voteEpoch.Load() && pingMsg.Role == int32(MasterRole) {
			masterNode := &raftNode{}
			err := masterNode.init(&nodeInitParams{
				id:      pingMsg.MyId,
				address: pingMsg.MasterAddress,
				role:    MasterRole,
			})
			if err != nil {
				return
			}
			c.myself.currentEpoch.Store(0)
			c.myself.lastCommitedEpoch.Store(0)
			msg := &prpc.RaftMsg{
				Type: prpc.RaftMsg_ReqActions,
				ReqActions: &prpc.RaftReqActions{
					MyId:          c.myself.id,
					CommitedEpoch: c.myself.lastCommitedEpoch.Load(),
				},
			}
			masterNode.send(msg)
			c.nodes[masterNode.id] = masterNode
			c.actors[masterNode.id] = masterNode
			c.master = masterNode
			if c.myself.role == MasterRole {
				c.myself.role = ActorRole
			}
			c.isMaster.Store(false)
			c.lastAckTime = time.Now()
		}
		return
	}

	if pingMsg.VoteEpoch > c.voteEpoch.Load() {
		if pingMsg.Role == int32(MasterRole) {
			var srcMasterId string
			if c.master != nil {
				srcMasterId = c.master.id
				c.master.role = ActorRole
			}
			fmt.Printf("change master:%s->%s\n", srcMasterId, node.id)
			node.role = MasterRole
			c.master = node
			c.isMaster.Store(c.myself == node)
			c.lastAckTime = time.Now()
		}
		c.voteEpoch.Store(pingMsg.VoteEpoch)
		c.nextVoteEpoch = pingMsg.VoteEpoch + 1
	}

	if node.lastCommitedEpoch.Load() != pingMsg.CommitedEpoch {
		node.changeCommitedEpoch(pingMsg.CommitedEpoch)
	}
	if node.currentEpoch.Load() != pingMsg.CurrentEpoch {
		node.changeCurrentEpoch(pingMsg.CurrentEpoch)
	}

	if c.isMaster.Load() {
		if node.currentEpoch.Load() < c.master.currentEpoch.Load() && time.Since(node.lastSyncActionsTime) > time.Second*2 {
			c.syncActionsLocked(node)
		}
	}

	for _, ns := range pingMsg.NodeStates {
		if ns.State == int32(MaybeDownNodeState) {
			fmt.Printf("from:%s mark target:%s as down\n", pingMsg.MyId, ns.MyId)
			_, okt := c.nodes[ns.MyId]
			if !okt {
				continue
			}
			node.updateOtherNodeState(ns.MyId, MaybeDownNodeState)
			if node.role == ActorRole {
				c.checkDownLocked(ns.MyId)
			}
		}
	}

	node.send(c.newPongMsg())
}

func (c *Cluster) handlePongMsgLock(msg *prpc.RaftMsg) {
	pong := msg.Pong
	if c.master != nil && pong.MyId == c.master.id {
		c.lastAckTime = time.Now()
	}
	st := c.myself.getOtherNodeState(pong.MyId)
	if st == MaybeDownNodeState || st == DownNodeState {
		c.myself.updateOtherNodeState(pong.MyId, ConnectedNodeState)
	}
}

func (c *Cluster) addNodeEpochLocked(nodeId string, epoch int64) {
	majority := len(c.actors)/2 + 1

	node, ok := c.nodes[nodeId]
	if !ok {
		return
	}
	node.changeCurrentEpoch(epoch)

	if epoch > c.myself.lastCommitedEpoch.Load() {
		count := 0
		for _, n := range c.actors {
			if n.currentEpoch.Load() >= epoch {
				count++
			}
		}
		fmt.Printf("check epoch:%d, count:%d, majority:%d, lastCommitedEpoch:%d\n", epoch, count, majority, c.myself.lastCommitedEpoch.Load())
		if count >= majority {
			c.handleCommitedActionLocked(epoch)
			pingMsg := c.newPingMsg()
			for _, n := range c.nodes {
				n.send(pingMsg)
			}
		}
	}
}

func (c *Cluster) handleCommitedActionLocked(commitedEpoch int64) {
	if commitedEpoch < c.myself.lastCommitedEpoch.Load() {
		return
	}
	if commitedEpoch > c.myself.currentEpoch.Load() {
		commitedEpoch = c.myself.currentEpoch.Load()
	}
	for e := c.myself.lastCommitedEpoch.Load() + 1; e <= commitedEpoch; e++ {
		if a, ok := c.actions[e]; ok {
			fmt.Printf("[action] commited action epoch:%d type:%d\n", a.Epoch, a.Type)
			if a.Type == prpc.RaftTransaction_NewNode {
				c.handleNewNodeActionLocked(a.NewNode)
			} else if a.Type == prpc.RaftTransaction_HashSlotAction {
				c.handleHashSlotActionLocked(a.HashSlot)
			}
		}
	}
	c.myself.changeCommitedEpoch(commitedEpoch)
	if c.master != nil {
		c.master.send(c.newPingMsg())
	} else {
		fmt.Printf("[action] not found master\n")
	}
}

func (c *Cluster) handleNewNodeActionLocked(action *prpc.NewNodeAction) {
	fmt.Printf("[action] new node:%v\n", action)
	if action.MyId == c.myself.id {
		return
	}

	fmt.Printf("[action] new node:%v\n", action)
	node := &raftNode{}
	node.init(&nodeInitParams{
		id:      action.MyId,
		address: action.MyAddress,
		role:    NodeRole(action.Role),
	})

	c.nodes[node.id] = node
	if node.role == ActorRole {
		c.actors[node.id] = node
	}
	if c.isMaster.Load() {
		c.syncActionsLocked(node)
	}
	node.send(c.newPingMsg())
}

func (c *Cluster) handleHashSlotActionLocked(action *prpc.HashSlotAction) {
	if len(action.Slots) != HashSlotSize/8 {
		return
	}
	node, ok := c.actors[action.MyId]
	if !ok {
		return
	}
	var formatSlots strings.Builder
	lastBegin := -1
	for i := 0; i < HashSlotSize; i++ {
		if !isSet(action.Slots, i) {
			if lastBegin >= 0 {
				formatSlots.WriteString(fmt.Sprintf(" %d->%d", lastBegin, i-1))
				lastBegin = -1
			}
		} else {
			if lastBegin == -1 {
				lastBegin = i
			}
			if i == HashSlotSize-1 {
				formatSlots.WriteString(fmt.Sprintf(" %d->%d", lastBegin, i))
			}
		}
	}
	fmt.Printf("[action] slot id:%s, slots:%s\n", action.MyId, formatSlots.String())

	if action.Step == SwitchSlotStep {
		for i := 0; i < HashSlotSize; i++ {
			if !isSet(action.Slots, i) {
				continue
			}
			node.setSlot(i)
			c.slots2Node[i] = node
			if c.myself != node && c.myself.isSet(i) {
				c.myself.unsetSlot(i)
			}
		}
	}
}

func (c *Cluster) markNodeMaybeDownLocked(targetId string) {
	targetNode, okt := c.nodes[targetId]
	if !okt || targetNode == c.myself {
		return
	}
	ls := c.myself.getOtherNodeState(targetId)
	if ls == DownNodeState {
		return
	}
	msg := c.newPingMsg()
	msg.Ping.NodeStates = append(msg.Ping.NodeStates, &prpc.NodeState{
		MyId:  targetId,
		State: int32(MaybeDownNodeState),
	})
	for _, n := range c.nodes {
		n.send(msg)
	}
	fmt.Printf("mark id:%s as MaybeDownNodeState\n", targetId)
	c.myself.updateOtherNodeState(targetId, MaybeDownNodeState)
	c.checkDownLocked(targetId)
}

func (c *Cluster) checkDownLocked(targetId string) {
	majority := len(c.actors)/2 + 1
	count := 0
	for _, n := range c.actors {
		if n.getOtherNodeState(targetId) == MaybeDownNodeState {
			count++
		}
	}
	if count >= majority {
		fmt.Printf("node:%s is down count:%d majority:%d\n", targetId, count, majority)
		c.markNodeisDownLocked(targetId)
	}
}

func (c *Cluster) markNodeisDownLocked(targetId string) {
	if c.myself.getOtherNodeState(targetId) == DownNodeState {
		fmt.Printf("targetId:%s is down", targetId)
		return
	}
	fmt.Printf("targetId:%s is down, master id:%s\n", targetId, c.master.id)
	c.myself.updateOtherNodeState(targetId, DownNodeState)
	if targetId == c.master.id {
		c.startElectionLocked()
	}
}

func (c *Cluster) getNextElectionDuration() time.Duration {
	rank := 0
	for _, n := range c.actors {
		state := c.myself.getOtherNodeState(n.id)
		if state == MaybeDownNodeState || state == DownNodeState {
			continue
		}
		if n.currentEpoch.Load() > c.myself.currentEpoch.Load() {
			rank++
		}
	}
	return time.Millisecond * time.Duration(1000+(c.rand.Int()%500)+rank*200)
}

func (c *Cluster) startElectionLocked() {
	if c.isMaster.Load() || c.myself.role != ActorRole || c.myself.getOtherNodeState(c.master.id) != DownNodeState {
		fmt.Printf("[election] exit election. master id:%s myrole:%d masterState:%d\n", c.master.id, c.myself.role, c.myself.getOtherNodeState(c.master.id))
		return
	}

	defer func() {
		c.wg.Add(1)
		go func() {
			defer c.wg.Add(-1)
			d := c.getNextElectionDuration()
			fmt.Printf("[election] set election time out:%d\n", d)
			t := time.NewTimer(d)
			select {
			case <-t.C:
				c.mtx.Lock()
				c.startElectionLocked()
				c.mtx.Unlock()
			case <-c.shutdownCtx.Done():
				return
			}
		}()
	}()

	if c.myself.lastVoteEpoch >= c.nextVoteEpoch {
		c.nextVoteEpoch = c.myself.lastVoteEpoch + 1
		return
	}
	nextVoteEpoch := c.nextVoteEpoch
	c.nextVoteEpoch++
	fmt.Printf("[election] start election my:%s epoch:%d\n", c.myself.id, nextVoteEpoch)
	c.state = ElectionCluster
	c.myself.lastVoteEpoch = nextVoteEpoch
	c.myself.voteId = c.myself.id
	electionMsg := prpc.RaftMsg{
		Type: prpc.RaftMsg_Election,
		Election: &prpc.RaftElection{
			MyId:      c.myself.id,
			VoteEpoch: nextVoteEpoch,
		},
	}
	for _, n := range c.actors {
		n.send(&electionMsg)
	}
}

func (c *Cluster) handleElectionMsgLocked(msg *prpc.RaftMsg) {
	election := msg.Election
	if c.myself.role != ActorRole || c.myself.lastVoteEpoch >= election.VoteEpoch {
		return
	}

	n, ok := c.actors[election.MyId]
	if !ok {
		return
	}
	if c.myself.currentEpoch.Load() > n.currentEpoch.Load() {
		return
	}

	c.myself.lastVoteEpoch = election.VoteEpoch
	c.myself.voteId = election.MyId

	electionRetMsg := prpc.RaftMsg{
		Type: prpc.RaftMsg_ElectionRet,
		ElectionRet: &prpc.RaftElectionRet{
			MyId:      c.myself.id,
			GotVoteId: n.id,
			VoteEpoch: election.VoteEpoch,
			Success:   true,
		},
	}
	n.send(&electionRetMsg)
}

func (c *Cluster) handleElectionRetMsgLocked(msg *prpc.RaftMsg) {
	er := msg.ElectionRet
	if !er.Success {
		return
	}
	hisNode, ok := c.actors[er.MyId]
	if !ok {
		return
	}
	fmt.Printf("[election] recv vote, from:%s obj:%s epoch:%d\n", er.MyId, er.GotVoteId, er.VoteEpoch)
	if hisNode.lastVoteEpoch < er.VoteEpoch {
		hisNode.lastVoteEpoch = er.VoteEpoch
		hisNode.voteId = er.GotVoteId

		if hisNode.voteId == c.myself.id && c.state == ElectionCluster && c.myself.voteId == c.myself.id {
			majority := len(c.actors)/2 + 1
			count := 0
			for _, n := range c.actors {
				if n.lastVoteEpoch == c.myself.lastVoteEpoch && n.voteId == c.myself.id {
					count++
				}
			}
			if count >= majority {
				c.becomeMasterLocked()
			}
		}
	}
}

func (c *Cluster) becomeMasterLocked() {
	fmt.Printf("myself:%s become master\n", c.myself.id)
	c.master = c.myself
	c.isMaster.Store(true)
	c.myself.role = MasterRole
	c.state = NormalCluster
	c.voteEpoch.Store(c.myself.lastVoteEpoch)
	c.nextVoteEpoch = c.myself.lastVoteEpoch + 1
	msg := c.newPingMsg()
	for _, n := range c.nodes {
		n.send(msg)
	}
}

func (c *Cluster) getRandOtherNodeLocked() *raftNode {
	l := len(c.nodes)
	if l == 1 {
		return nil
	}
	for {
		r := c.rand.Int() % l
		i := 0
		for _, node := range c.nodes {
			if i == r && node != c.myself {
				return node
			}
			i++
		}
	}
}
