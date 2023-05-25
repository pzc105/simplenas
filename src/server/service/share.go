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
	UserId user.ID
	ItemId category.ID
}

type ShareInfo struct {
	ShareId     string
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

type ShareManager struct {
	mtx     sync.Mutex
	shares  map[string]*ShareInfo
	timeOut timeHeap

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
			heap.Push(&sm.timeOut, &s)
		}
	}
	sm.nextId.Store(0)
}

type ShareCategoryItemParams struct {
	UserId    user.ID
	ItemId    category.ID
	MaxCount  int
	ExpiresAt time.Time
}

func (sm *ShareManager) ShareCategoryItem(params *ShareCategoryItemParams) (shareid string, err error) {
	si := &ShareInfo{
		ShareId:   sm.genShareId(),
		MaxCount:  params.MaxCount,
		ExpiresAt: params.ExpiresAt,

		ShareItemInfo: &ShareItemInfo{
			UserId: params.UserId,
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
			osi := sm.timeOut[0]
			delete(sm.shares, osi.ShareId)
			heap.Pop(&sm.timeOut)
		} else {
			break
		}
	}
	sm.shares[si.ShareId] = si
	heap.Push(&sm.timeOut, si)
	return si.ShareId, nil
}

func (sm *ShareManager) GetShareItemInfo(shareid string) (*ShareItemInfo, error) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	si, ok := sm.shares[shareid]
	if ok && si.ShareItemInfo != nil && si.ExpiresAt.Before(time.Now()) {
		return si.ShareItemInfo, nil
	}
	return nil, errors.New("not found")
}
