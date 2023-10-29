package service

import (
	"container/heap"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pnas/category"
	"pnas/db"
	"pnas/user"
	"pnas/utils"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slices"
)

const (
	sharePrefixRedisKey      = "share_"
	shareCountPrefixRedisKey = "share_count_"
)

var ShareIdPrefix string

func init() {
	ShareIdPrefix = "P"
}

type ShareItemInfo struct {
	ItemId category.ID
}

type ShareInfo struct {
	ShareId     string
	UserId      user.ID
	UseCounting bool
	MaxCount    int
	ExpiresAt   time.Time

	ShareItemInfo *ShareItemInfo
}

type timeHeap []*ShareInfo

func (h timeHeap) Len() int           { return len(h) }
func (h timeHeap) Less(i, j int) bool { return h[i].ExpiresAt.Before(h[j].ExpiresAt) }
func (h timeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *timeHeap) Push(x any)        { *h = append(*h, x.(*ShareInfo)) }
func (h *timeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type SharesInterface interface {
	ShareCategoryItem(params *ShareCategoryItemParams) (shareid string, err error)
	DelShare(shareid string) error
	GetShareItemInfo(shareid string) (*ShareInfo, error)
	GetUserSharedItemInfos(userId user.ID) []*ShareInfo
}

type ShareManager struct {
	SharesInterface
	mtx        sync.Mutex
	shares     map[string]*ShareInfo
	userShares map[user.ID][]string
	timeOut    timeHeap

	nextId atomic.Int64
}

func shareObjectRedisKey(shareId string) string {
	return sharePrefixRedisKey + shareId
}

func (sm *ShareManager) genShareId() string {
	nid := utils.FetchAndAdd(&sm.nextId, int64(1))
	return fmt.Sprintf("%s%d%d", ShareIdPrefix, time.Now().UnixMicro(), nid)
}

func (sm *ShareManager) Init() {
	sm.shares = make(map[string]*ShareInfo)
	sm.userShares = make(map[user.ID][]string)
	keys, err := db.GREDIS.Keys(context.Background(), shareObjectRedisKey("*")).Result()
	if err == nil {
		for _, k := range keys {
			objectStr, err := db.GREDIS.Get(context.Background(), k).Result()
			if err != nil {
				continue
			}
			var s ShareInfo
			json.Unmarshal([]byte(objectStr), &s)
			sm.shares[s.ShareId] = &s
			sm.userShares[s.UserId] = append(sm.userShares[s.UserId], s.ShareId)
			heap.Push(&sm.timeOut, &s)
		}
	}
	sm.nextId.Store(0)
}

type ShareCategoryItemParams struct {
	UserId    user.ID
	ItemOwner user.ID
	ItemId    category.ID
	MaxCount  int
	ExpiresAt time.Time
}

func (sm *ShareManager) _deleteShare(si *ShareInfo) {
	delete(sm.shares, si.ShareId)
	if si.ShareItemInfo != nil {
		i := slices.Index(sm.userShares[si.UserId], si.ShareId)
		if i != -1 {
			sm.userShares[si.UserId] = append(sm.userShares[si.UserId][:i],
				sm.userShares[si.UserId][i+1:]...)
		}
	}
	db.GREDIS.Del(context.Background(), shareObjectRedisKey(si.ShareId))
}

func (sm *ShareManager) ShareCategoryItem(params *ShareCategoryItemParams) (shareid string, err error) {
	si := &ShareInfo{
		ShareId:   sm.genShareId(),
		UserId:    params.UserId,
		MaxCount:  params.MaxCount,
		ExpiresAt: params.ExpiresAt,

		ShareItemInfo: &ShareItemInfo{
			ItemId: params.ItemId,
		},
	}
	if params.MaxCount > 0 {
		si.UseCounting = true
	}
	objectStr, _ := json.Marshal(si)
	db.GREDIS.Set(context.Background(), shareObjectRedisKey(si.ShareId), objectStr, time.Until(si.ExpiresAt))

	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	for len(sm.timeOut) > 0 {
		if sm.timeOut[0].ExpiresAt.Before(time.Now()) {
			sm._deleteShare(sm.timeOut[0])
			heap.Pop(&sm.timeOut)
		} else {
			break
		}
	}
	sm.shares[si.ShareId] = si
	sm.userShares[params.UserId] = append(sm.userShares[params.UserId], si.ShareId)
	heap.Push(&sm.timeOut, si)
	return si.ShareId, nil
}

func (sm *ShareManager) DelShare(shareid string) error {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	si, ok := sm.shares[shareid]
	if !ok {
		return errors.New("not found")
	}
	sm._deleteShare(si)
	return nil
}

func (sm *ShareManager) GetShareItemInfo(shareid string) (*ShareInfo, error) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	si, ok := sm.shares[shareid]
	if ok && si.ShareItemInfo != nil && si.ExpiresAt.After(time.Now()) {
		return si, nil
	}
	return nil, errors.New("not found")
}

func (sm *ShareManager) GetUserSharedIds(userId user.ID) []string {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	ss, ok := sm.userShares[userId]
	if !ok {
		return []string{}
	}
	ret := make([]string, len(ss))
	copy(ret, ss)
	return ret
}

func (sm *ShareManager) GetUserSharedItemInfos(userId user.ID) []*ShareInfo {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	ss, ok := sm.userShares[userId]
	if !ok {
		return []*ShareInfo{}
	}
	var ret []*ShareInfo
	for _, shareid := range ss {
		si, ok := sm.shares[shareid]
		if !ok || si.ShareItemInfo == nil || si.ExpiresAt.Before(time.Now()) {
			continue
		}
		ret = append(ret, si)
	}
	return ret
}
