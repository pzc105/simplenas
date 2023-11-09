package category

import (
	"fmt"
	"pnas/db"
	"pnas/log"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Manager struct {
	itemsMtx sync.Mutex
	items    map[ID]*CategoryItem

	dbMapMtx    sync.Mutex
	dbItemMtxes map[ID]*sync.Mutex
}

func (m *Manager) Init() {
	m.items = make(map[ID]*CategoryItem)
	m.dbItemMtxes = make(map[ID]*sync.Mutex)
}

func (m *Manager) requireDbMtx(itemId ID) *sync.Mutex {
	m.dbMapMtx.Lock()
	defer m.dbMapMtx.Unlock()
	dbmtx, ok := m.dbItemMtxes[itemId]
	if !ok {
		dbmtx = &sync.Mutex{}
		m.dbItemMtxes[itemId] = dbmtx
	}
	return dbmtx
}

func (m *Manager) delDbMtx(itemId ID) {
	m.dbMapMtx.Lock()
	defer m.dbMapMtx.Unlock()
	_, ok := m.dbItemMtxes[itemId]
	if ok {
		delete(m.dbItemMtxes, itemId)
	}
}

func (m *Manager) addItem(item *CategoryItem) {
	m.itemsMtx.Lock()
	defer m.itemsMtx.Unlock()
	_, ok := m.items[item.base.Id]
	if ok {
		log.Warnf("[category] duplicate add item: %d", item.base.Id)
		return
	}
	m.items[item.base.Id] = item
}

func (m *Manager) queryItem(itemId ID) *CategoryItem {
	m.itemsMtx.Lock()
	defer m.itemsMtx.Unlock()
	item, ok := m.items[itemId]
	if ok {
		return item
	}
	return nil
}

func (m *Manager) removeItem(itemId ID) {
	m.itemsMtx.Lock()
	defer m.itemsMtx.Unlock()
	_, ok := m.items[itemId]
	if ok {
		delete(m.items, itemId)
	}
}

func (m *Manager) AddItem(params *NewCategoryParams) (*CategoryItem, error) {
	querier := params.Creator
	if params.Sudo {
		querier = AdminId
	}
	parentItem, err := m.GetItem(querier, params.ParentId)
	if err != nil {
		return nil, err
	}
	if !parentItem.IsDirectory() {
		return nil, errors.New("isn't a directory")
	}
	if !parentItem.HasWriteAuth(params.Creator) && !params.Sudo {
		return nil, errors.New("not auth")
	}

	dbmtx := m.requireDbMtx(params.ParentId)
	dbmtx.Lock()
	defer dbmtx.Unlock()

	item, err := addItem(params)
	if err == nil {
		parentItem.addedSubItem(item.base.Id)
		m.addItem(item)
	}
	return item, err
}

func (m *Manager) GetItem(querier int64, itemId ID) (*CategoryItem, error) {
	item := m.queryItem(itemId)
	if item != nil {
		if !item.HasReadAuth(querier) {
			return nil, errors.New("no auth")
		}
		return item, nil
	}

	dbmtx := m.requireDbMtx(itemId)
	dbmtx.Lock()
	defer dbmtx.Unlock()

	item = m.queryItem(itemId)
	if item != nil {
		return item, nil
	}

	item, err := _loadItem(itemId)
	if err != nil {
		m.delDbMtx(itemId)
		return item, err
	}

	m.itemsMtx.Lock()
	m.items[itemId] = item
	m.itemsMtx.Unlock()
	if !item.HasReadAuth(querier) {
		return nil, errors.New("no auth")
	}
	return item, err
}

func (m *Manager) GetItemByName(querier int64, parentId ID, name string) (*CategoryItem, error) {
	if parentId <= 0 {
		return nil, errors.New("wrong parent id")
	}
	id, err := _loadItemIdByName(parentId, name)
	if err != nil {
		return nil, err
	}
	return m.GetItem(querier, id)
}

func (m *Manager) GetItemsByParent(params *GetItemsByParentParams) ([]*CategoryItem, error) {
	item, err := m.GetItem(params.Querier, params.ParentId)
	if err != nil {
		return nil, err
	}
	if params.PageNum >= 0 && params.Rows > 0 {
		subIds := item.GetSubItemIds()
		offset := params.PageNum * params.Rows
		if int(offset) >= len(subIds) {
			return nil, errors.New("out of range")
		}
		rows := params.Rows
		if rows > int32(len(subIds))-offset {
			rows = int32(len(subIds)) - offset
		}
		ids := subIds[offset : offset+rows]
		return m.GetItems(params.Querier, ids...)
	}
	return m.GetItems(params.Querier, item.subItemIds...)
}

func (m *Manager) GetItems(querier int64, itemIds ...ID) ([]*CategoryItem, error) {
	remainIds := make([]ID, 0, len(itemIds))
	ret := make([]*CategoryItem, 0, len(itemIds))
	m.itemsMtx.Lock()
	for _, id := range itemIds {
		item, ok := m.items[id]
		if ok {
			ret = append(ret, item)
		} else {
			remainIds = append(remainIds, id)
		}
	}
	m.itemsMtx.Unlock()
	if len(remainIds) == 0 {
		return ret, nil
	}

	mtxes := make([]*sync.Mutex, 0, len(remainIds))
	realNeedQueryIds := make([]ID, 0, len(remainIds))
	for _, itemId := range remainIds {
		mtx := m.requireDbMtx(itemId)
		mtxes = append(mtxes, mtx)
		mtx.Lock()
		item := m.queryItem(itemId)
		if item != nil && (item.HasReadAuth(querier)) {
			ret = append(ret, item)
		} else {
			realNeedQueryIds = append(realNeedQueryIds, itemId)
		}
	}

	defer func() {
		for _, mtx := range mtxes {
			mtx.Unlock()
		}
	}()

	items, err := _loadItems(realNeedQueryIds...)
	if err == nil {
		for _, item := range items {
			if item.HasReadAuth(querier) {
				ret = append(ret, item)
			}
		}
	}

	return ret, err
}

func (m *Manager) DelItem(deleter int64, itemId ID) (err error) {
	item, err := m.GetItem(deleter, itemId)
	if err != nil {
		return err
	}
	if !item.HasWriteAuth(deleter) {
		return errors.New("not auth")
	}

	var lookingItems []*CategoryItem
	parentItem, _ := m.GetItem(AdminId, item.base.ParentId)
	parentDbMtx := m.requireDbMtx(item.base.ParentId)
	lookingItems = append(lookingItems, item)

	parentDbMtx.Lock()
	defer parentDbMtx.Unlock()

	toDelItemsDbMtx := []*sync.Mutex{m.requireDbMtx(itemId)}
	toDelItemsDbMtx[0].Lock()
	defer func() {
		if err != nil {
			for _, mtx := range toDelItemsDbMtx {
				mtx.Unlock()
			}
		}
	}()

	lockedIds := []ID{itemId}

	for len(lookingItems) > 0 {
		item := lookingItems[len(lookingItems)-1]
		lookingItems = lookingItems[:len(lookingItems)-1]
		items, _ := m.GetItems(AdminId, item.GetSubItemIds()...)
		lookingItems = append(lookingItems, items...)
		for _, item := range items {
			dbmtx := m.requireDbMtx(item.base.Id)
			dbmtx.Lock()
			toDelItemsDbMtx = append(toDelItemsDbMtx, dbmtx)
			lockedIds = append(lockedIds, item.base.Id)
		}
	}

	var sb strings.Builder
	for _, id := range lockedIds {
		if sb.Len() == 0 {
			sb.WriteString(fmt.Sprint(id))
		} else {
			sb.WriteString("," + fmt.Sprint(id))
		}
	}

	sql := fmt.Sprintf("delete from category_items where id in (%s)", sb.String())
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	for _, id := range lockedIds {
		m.removeItem(id)
	}
	for _, mtx := range toDelItemsDbMtx {
		mtx.Unlock()

	}
	for _, id := range lockedIds {
		m.delDbMtx(id)
	}

	if parentItem != nil {
		parentItem.deletedSubItem(itemId)
	}
	return nil
}

func (m *Manager) IsRelationOf(itemId ID, parentId ID) bool {
	_, err := m.GetItem(AdminId, parentId)
	if err != nil {
		return false
	}
	var nextParentId = itemId
	for {
		item, err := m.GetItem(AdminId, nextParentId)
		if err != nil {
			return false
		}
		ii := item.GetItemInfo()
		if ii.Id == parentId {
			return true
		}
		nextParentId = ii.ParentId
		if nextParentId <= 0 {
			return false
		}
	}
}

type SearchParams struct {
	Querier      int64
	RootId       ID
	ExistedWords string
	PageNum      int32
	Rows         int32
}

func (m *Manager) SearchRows(params *SearchParams) (int, error) {
	if params.PageNum < 0 || params.Rows <= 0 {
		return -1, errors.New("invalid params")
	}
	sql := fmt.Sprintf("select id from category_items where match (name, introduce) against ('%s')", params.ExistedWords)
	rows, err := db.Query(sql)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	ret := 0
	for rows.Next() {
		var id ID
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		if params.RootId <= 0 || m.IsRelationOf(id, params.RootId) {
			ret += 1
		}
	}
	return ret, nil
}

func (m *Manager) Search(params *SearchParams) ([]*CategoryItem, error) {
	if params.PageNum < 0 || params.Rows <= 0 {
		return nil, errors.New("invalid params")
	}
	sql := fmt.Sprintf("select id from category_items where match (name, introduce) against ('%s') limit ?, ?", params.ExistedWords)
	rows, err := db.Query(sql, params.PageNum*params.Rows, params.Rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []*CategoryItem
	for rows.Next() {
		var id ID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		if params.RootId <= 0 || m.IsRelationOf(id, params.RootId) {
			item, err := m.GetItem(params.Querier, id)
			if err != nil {
				log.Warnf("[category] %v", err)
				continue
			}
			ret = append(ret, item)
		}
	}
	return ret, nil
}
