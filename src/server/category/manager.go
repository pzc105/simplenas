package category

import (
	"pnas/db"
	"pnas/log"
	"sync"
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

func (m *Manager) NewItem(params *NewCategoryParams) (*CategoryItem, error) {
	parentItem, err := m.GetItem(params.ParentId)
	if err != nil {
		return nil, err
	}

	dbmtx := m.requireDbMtx(params.ParentId)
	dbmtx.Lock()
	defer dbmtx.Unlock()

	item, err := newItem(params)
	if err == nil {
		parentItem.addedSubItem(item.base.Id)
	}
	m.addItem(item)
	return item, err
}

func (m *Manager) GetItem(itemId ID) (*CategoryItem, error) {
	item := m.queryItem(itemId)
	if item != nil {
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

	return item, err
}

func (m *Manager) GetItems(itemIds ...ID) ([]*CategoryItem, error) {
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
		if item != nil {
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
		ret = append(ret, items...)
	}

	return ret, err
}

func (m *Manager) DelItem(itemId ID) error {
	item, err := m.GetItem(itemId)
	if err != nil {
		return err
	}

	var toDelItems []*CategoryItem
	parentItem, _ := m.GetItem(item.base.ParentId)
	parentDbMtx := m.requireDbMtx(item.base.ParentId)
	toDelItems = append(toDelItems, item)
	toDelItemsDbMtx := []*sync.Mutex{m.requireDbMtx(itemId)}
	parentDbMtx.Lock()
	defer parentDbMtx.Unlock()
	toDelItemsDbMtx[0].Lock()

	flags := make(map[ID]bool)

	for len(toDelItems) > 0 {
		item := toDelItems[len(toDelItems)-1]
		if _, ok := flags[item.base.Id]; !ok {
			flags[item.base.Id] = true
			items, _ := m.GetItems(item.GetSubItemIds()...)
			toDelItems = append(toDelItems, items...)
			for _, item := range items {
				dbmtx := m.requireDbMtx(item.base.Id)
				toDelItemsDbMtx = append(toDelItemsDbMtx, dbmtx)
				dbmtx.Lock()
			}
		} else {
			dbmtx := toDelItemsDbMtx[len(toDelItemsDbMtx)-1]
			toDelItems = toDelItems[:len(toDelItems)-1]
			toDelItemsDbMtx = toDelItemsDbMtx[:len(toDelItemsDbMtx)-1]
			_, err = db.Exec("call del_category(?)", item.base.Id)
			dbmtx.Unlock()
			if err != nil {
				log.Warnf("[category] id:%d del error: %v", item.base.Id, err)
				continue
			}
			m.removeItem(item.base.Id)
			m.delDbMtx(item.base.Id)
		}
	}
	if parentItem != nil {
		parentItem.deletedSubItem(itemId)
	}
	if len(toDelItemsDbMtx) > 0 {
		log.Panic("del item")
	}
	return nil
}
