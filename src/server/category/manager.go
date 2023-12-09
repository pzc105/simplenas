package category

import (
	"database/sql"
	"fmt"
	"pnas/db"
	"pnas/log"
	"pnas/ptype"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Manager struct {
	itemsMtx sync.Mutex
	// TODO: support LFU
	items map[ptype.CategoryID]*CategoryItem

	dbMapMtx    sync.Mutex
	dbItemMtxes map[ptype.CategoryID]*sync.Mutex
}

func (m *Manager) Init() {
	m.items = make(map[ptype.CategoryID]*CategoryItem)
	m.dbItemMtxes = make(map[ptype.CategoryID]*sync.Mutex)
}

func (m *Manager) requireDbMtx(itemId ptype.CategoryID) *sync.Mutex {
	m.dbMapMtx.Lock()
	defer m.dbMapMtx.Unlock()
	dbmtx, ok := m.dbItemMtxes[itemId]
	if !ok {
		dbmtx = &sync.Mutex{}
		m.dbItemMtxes[itemId] = dbmtx
	}
	return dbmtx
}

func (m *Manager) delDbMtx(itemId ptype.CategoryID) {
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

func (m *Manager) queryItem(itemId ptype.CategoryID) *CategoryItem {
	m.itemsMtx.Lock()
	defer m.itemsMtx.Unlock()
	item, ok := m.items[itemId]
	if ok {
		return item
	}
	return nil
}

func (m *Manager) removeItem(itemId ptype.CategoryID) {
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
		querier = ptype.AdminId
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
		log.Debugf("[category] user %d add item %d type: %d name: %s", params.Creator, item.base.Id, params.TypeId, params.Name)
		item.base.Other = params.Other
		parentItem.addedSubItem(item.base.Id)
		m.addItem(item)
	}
	return item, err
}

func (m *Manager) GetItem(querier ptype.UserID, itemId ptype.CategoryID) (*CategoryItem, error) {
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

func (m *Manager) RefreshItem(itemId ptype.CategoryID) error {
	dbmtx := m.requireDbMtx(itemId)
	dbmtx.Lock()
	defer dbmtx.Unlock()
	item, err := _loadItem(itemId)
	if err != nil {
		m.delDbMtx(itemId)
		return err
	}
	m.itemsMtx.Lock()
	m.items[itemId] = item
	m.itemsMtx.Unlock()
	return nil
}

func (m *Manager) GetItemByName(querier ptype.UserID, parentId ptype.CategoryID, name string) (*CategoryItem, error) {
	if parentId <= 0 {
		return nil, errors.New("wrong parent id")
	}
	id, err := _loadItemIdByName(parentId, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
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
		if len(subIds) == 0 {
			return []*CategoryItem{}, nil
		}
		offset := params.PageNum * params.Rows
		if int(offset) >= len(subIds) {
			return nil, errors.New("out of range")
		}
		rows := params.Rows
		if rows > int32(len(subIds))-offset {
			rows = int32(len(subIds)) - offset
		}
		var ids []ptype.CategoryID
		if params.Desc {
			for i := len(subIds) - 1 - int(offset); i >= 0; i-- {
				ids = append(ids, subIds[i])
				if len(ids) >= int(rows) {
					break
				}
			}
		} else {
			ids = subIds[offset : offset+rows]
		}
		return m.GetItems(params.Querier, ids...)
	}
	return m.GetItems(params.Querier, item.subItemIds...)
}

func (m *Manager) GetItems(querier ptype.UserID, itemIds ...ptype.CategoryID) ([]*CategoryItem, error) {
	remainIds := make([]ptype.CategoryID, 0, len(itemIds))
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
	realNeedQueryIds := make([]ptype.CategoryID, 0, len(remainIds))
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
	m.itemsMtx.Lock()
	for _, item := range items {
		m.items[item.base.Id] = item
	}
	m.itemsMtx.Unlock()
	if err == nil {
		for _, item := range items {
			if item.HasReadAuth(querier) {
				ret = append(ret, item)
			}
		}
	}

	return ret, err
}

func (m *Manager) DelItem(deleter ptype.UserID, itemId ptype.CategoryID) (err error) {
	item, err := m.GetItem(deleter, itemId)
	if err != nil {
		return err
	}
	if !item.HasWriteAuth(deleter) {
		return errors.New("not auth")
	}

	var lookingItems []*CategoryItem
	parentItem, _ := m.GetItem(ptype.AdminId, item.base.ParentId)
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

	lockedIds := []ptype.CategoryID{itemId}

	for len(lookingItems) > 0 {
		item := lookingItems[len(lookingItems)-1]
		lookingItems = lookingItems[:len(lookingItems)-1]
		items, _ := m.GetItems(ptype.AdminId, item.GetSubItemIds()...)
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

	log.Debugf("[category] user %d del item %s", deleter, sb.String())

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

func (m *Manager) IsRelationOf(itemId ptype.CategoryID, parentId ptype.CategoryID) bool {
	_, err := m.GetItem(ptype.AdminId, parentId)
	if err != nil {
		return false
	}
	var nextParentId = itemId
	for {
		item, err := m.GetItem(ptype.AdminId, nextParentId)
		if err != nil {
			return false
		}
		ii := item.GetItemBaseInfo()
		if ii.Id == parentId {
			return true
		}
		nextParentId = ii.ParentId
		if nextParentId <= 0 {
			return false
		}
	}
}

func (m *Manager) SearchRows(params *SearchParams) (int, error) {
	if params.PageNum < 0 || params.Rows <= 0 {
		return -1, errors.New("invalid params")
	}
	var condBuild strings.Builder
	for _, word := range params.ExistedWords {
		if len(word) == 0 {
			continue
		}
		if condBuild.Len() > 0 {
			condBuild.WriteString(" ")
		}
		condBuild.WriteString("+")
		condBuild.WriteString(word)
	}
	for _, word := range params.NotExistedWords {
		if condBuild.Len() > 0 {
			condBuild.WriteString(" ")
		}
		condBuild.WriteString("-")
		condBuild.WriteString(word)
	}
	if condBuild.Len() == 0 {
		return -1, errors.New("invalid params")
	}
	sql := fmt.Sprintf("select id from category_items where match (name, introduce) against ('%s' in boolean mode)",
		condBuild.String())
	rows, err := db.Query(sql)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	ret := 0
	ids := []ptype.CategoryID{}
	for rows.Next() {
		var id ptype.CategoryID
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		ids = append(ids, id)
	}
	m.GetItems(ptype.AdminId, ids...)
	for _, id := range ids {
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
	var condBuild strings.Builder
	for _, word := range params.ExistedWords {
		if len(word) == 0 {
			continue
		}
		if condBuild.Len() > 0 {
			condBuild.WriteString(" ")
		}
		condBuild.WriteString("+")
		condBuild.WriteString(word)
	}
	for _, word := range params.NotExistedWords {
		if condBuild.Len() > 0 {
			condBuild.WriteString(" ")
		}
		condBuild.WriteString("-")
		condBuild.WriteString(word)
	}
	if condBuild.Len() == 0 {
		return nil, errors.New("invalid params")
	}
	sql := fmt.Sprintf("select id from category_items where match (name, introduce) against ('%s' in boolean mode)",
		condBuild.String())
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []*CategoryItem
	ids := []ptype.CategoryID{}
	for rows.Next() {
		var id ptype.CategoryID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	offset := int32(0)
	for _, id := range ids {
		if params.RootId <= 0 || m.IsRelationOf(id, params.RootId) {
			offset += 1
			if offset <= params.PageNum*params.Rows {
				continue
			}
			item, err := m.GetItem(params.Querier, id)
			if err != nil {
				log.Warnf("[category] %v", err)
				continue
			}
			ret = append(ret, item)
			if len(ret) >= int(params.Rows) {
				break
			}
		}
	}
	return ret, nil
}
