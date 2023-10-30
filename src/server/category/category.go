package category

import (
	"fmt"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/utils"
	"pnas/video"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type ID int64

const (
	RootId     = 1
	NotExisted = -2
)

const (
	AdminId = 1
)

type BaseItem struct {
	Id           ID
	Creator      int64
	CreatedAt    time.Time
	TypeId       prpc.CategoryItem_Type
	Name         string
	ResourcePath string
	PosterPath   string
	Introduce    string
	UpdatedAt    time.Time
	ParentId     ID
}

type CategoryItem struct {
	mtx        sync.Mutex
	base       BaseItem
	auth       utils.AuthBitSet
	subItemIds []ID
}

func _initSubItemIds(item *CategoryItem) error {
	sql := `select id from pnas.category_items where parent_id=?`
	rows, err := db.Query(sql, item.base.Id)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()
	var subIds []ID
	for rows.Next() {
		var itemId ID
		err := rows.Scan(&itemId)
		if err != nil {
			log.Warn(err)
			continue
		}
		subIds = append(subIds, itemId)
	}
	item.subItemIds = subIds
	return nil
}

func _loadItem(itemId ID) (*CategoryItem, error) {
	var item CategoryItem
	item.base.Id = itemId
	var byteAuth []byte
	sql := `select type_id, name, creator, auth, resource_path, poster_path, introduce, created_at, updated_at, parent_id
				from pnas.category_items
				where id=?`
	err := db.QueryRow(sql, itemId).Scan(
		&item.base.TypeId, &item.base.Name, &item.base.Creator, &byteAuth, &item.base.ResourcePath, &item.base.PosterPath,
		&item.base.Introduce, &item.base.CreatedAt, &item.base.UpdatedAt, &item.base.ParentId,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	item.auth = utils.NewBitSet(AuthMax)
	item.auth.UnmarshalBinary(byteAuth)

	err = _initSubItemIds(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func _loadItemIdByName(parentId ID, name string) (ID, error) {
	var ret ID
	sql := `select id from pnas.category_items where parent_id=? and name=?`
	err := db.QueryRow(sql, parentId, name).Scan(&ret)
	return ret, err
}

func _loadItems(itemIds ...ID) ([]*CategoryItem, error) {
	if len(itemIds) == 0 {
		return []*CategoryItem{}, nil
	}
	var conds []string
	for _, id := range itemIds {
		conds = append(conds, fmt.Sprintf("id=%d", id))
	}
	cond := strings.Join(conds, " or ")
	sql := `select id, type_id, name, creator, auth, resource_path, poster_path, introduce, created_at, updated_at, parent_id
					from pnas.category_items where ` + cond
	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	var items []*CategoryItem
	var byteAuth []byte
	for rows.Next() {
		var item CategoryItem
		err = rows.Scan(
			&item.base.Id, &item.base.TypeId, &item.base.Name, &item.base.Creator, &byteAuth,
			&item.base.ResourcePath, &item.base.PosterPath,
			&item.base.Introduce, &item.base.CreatedAt, &item.base.UpdatedAt, &item.base.ParentId,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		err = _initSubItemIds(&item)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		items = append(items, &item)
	}
	return items, nil
}

type NewCategoryParams struct {
	ParentId     ID
	Creator      int64
	TypeId       prpc.CategoryItem_Type
	Name         string
	ResourcePath string
	PosterPath   string
	Introduce    string
	Auth         utils.AuthBitSet
	CompareName  bool
}

func addItem(params *NewCategoryParams) (*CategoryItem, error) {
	var newId ID

	if params.CompareName {
		var c int
		err := db.QueryRow("select count(*) from pnas.category_items where parent_id=? and name=?", params.ParentId, params.Name).Scan(&c)
		if err != nil {
			return nil, err
		}
		if c > 0 {
			return nil, errors.New(fmt.Sprintf("existed name: %s", params.Name))
		}
	}

	if params.Auth.BitSet == nil {
		params.Auth = utils.NewBitSet(AuthMax)
	}
	byteAuth, err := params.Auth.MarshalBinary()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = db.QueryRow("call new_category(?, ?, ?, ?, ?, ?, ?, ?, @new_item_id)",
		params.TypeId,
		params.Name,
		params.Creator,
		byteAuth,
		params.ResourcePath,
		params.PosterPath,
		params.Introduce,
		params.ParentId).Scan(&newId)
	if err != nil {
		log.Warn(err)
		return nil, errors.WithStack(err)
	}
	if newId == NotExisted {
		return nil, errors.New("parent not existed")
	}
	return _loadItem(newId)
}

func (c *CategoryItem) addedSubItem(subItemId ID) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if slices.IndexFunc(c.subItemIds, func(sid ID) bool { return subItemId == sid }) != -1 {
		log.Warnf("[category] duplicate added sub id: %d", subItemId)
		return
	}
	c.subItemIds = append(c.subItemIds, subItemId)
}

func (c *CategoryItem) deletedSubItem(subItemId ID) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	index := slices.IndexFunc(c.subItemIds, func(sid ID) bool { return subItemId == sid })
	if index == -1 {
		log.Warnf("[category] not found sub id: %d", subItemId)
		return
	}
	c.subItemIds = append(c.subItemIds[:index], c.subItemIds[index+1:]...)
}

func (c *CategoryItem) GetItemInfo() BaseItem {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base
}

func (c *CategoryItem) Rename(newName string) error {
	sql := "update pnas.category_items set name=? where id=?"
	_, err := db.Exec(sql, newName, c.base.Id)
	if err != nil {
		return errors.WithStack(err)
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.base.Name = newName
	c.base.UpdatedAt = time.Now()
	return nil
}

func (c *CategoryItem) GetOwner() int64 {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.Creator
}

func (c *CategoryItem) HasAndAuths(auths ...uint) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, i := range auths {
		if !c.auth.Test(i) {
			return false
		}
	}
	return true
}

func (c *CategoryItem) HasReadAuth(who int64) bool {
	if who == AdminId {
		return true
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.base.Creator == who || c.auth.Test(AuthOtherRead) {
		return true
	}
	return false
}

func (c *CategoryItem) HasWriteAuth(who int64) bool {
	if who == AdminId {
		return true
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.base.Creator == who || c.auth.Test(AuthOtherWrite) {
		return true
	}
	return false
}

func (c *CategoryItem) GetType() prpc.CategoryItem_Type {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.TypeId
}

func (c *CategoryItem) IsDirectory() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.TypeId == prpc.CategoryItem_Directory || c.base.TypeId == prpc.CategoryItem_Home
}

func (c *CategoryItem) GetSubItemIds() []ID {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	ret := make([]ID, len(c.subItemIds))
	copy(ret, c.subItemIds)
	return ret
}

func (c *CategoryItem) UpdatePosterPath(path string) error {
	sql := "update pnas.category_items set poster_path=? where id=?"
	_, err := db.Exec(sql, path, c.base.Id)
	if err == nil {
		c.mtx.Lock()
		defer c.mtx.Unlock()
		c.base.PosterPath = path
	} else {
		log.Warnf("itemId: %d, %v", c.base.Id, err)
	}
	return err
}

func (c *CategoryItem) GetPosterPath() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.PosterPath
}

func (c *CategoryItem) GetVideoId() video.ID {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.base.TypeId != prpc.CategoryItem_Video {
		return -1
	}
	vid, err := strconv.ParseInt(c.base.ResourcePath, 10, 64)
	if err != nil {
		return -1
	}
	return video.ID(vid)
}
